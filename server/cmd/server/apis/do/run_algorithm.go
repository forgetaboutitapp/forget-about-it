package do

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/github.com/forgetaboutitapp/protobufs/scheduler"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"strings"
	"time"
)

type AlgorithmStruct struct {
	Alloc         string `json:"alloc"`
	ApiVersion    int    `json:"api-version"`
	Author        string `json:"author"`
	Dealloc       string `json:"dealloc"`
	Desc          string `json:"desc"`
	DownloadUrl   string `json:"download-url"`
	Init          string `json:"init"`
	License       string `json:"license"`
	ModuleName    string `json:"module-name"`
	AlgorithmName string `json:"algorithm-name"`
	RemoteURL     string `json:"remote-url"`
	Version       int    `json:"version"`
	WasmString    string `json:"wasm"`
	WasmBytes     []byte
}

type RunAlgorithm struct {
	algo           AlgorithmStruct
	allGrades      []sql_queries.QuestionsLog
	tagsByQuestion map[uint32][]string
	tagsToAsk      []string
}

func runAlgorithm(ctx context.Context, arg RunAlgorithm) (*scheduler.ResponseCorrect, error, string) {
	slog.Info("starting algorithm", "arg.allGrades", arg.allGrades)
	v := sha256.Sum256(arg.algo.WasmBytes)
	slog.Info("sha256sum", "sha", hex.EncodeToString(v[:]))
	runtime := wazero.NewRuntime(ctx)

	defer func(runtime wazero.Runtime, ctx context.Context) {
		err := runtime.Close(ctx)
		if err != nil {
			slog.Error("error closing runtime", "err", err)
		}
	}(runtime, ctx) // This closes everything this Runtime created.
	_, err := runtime.NewHostModuleBuilder(arg.algo.ModuleName).Instantiate(ctx)
	if err != nil {
		slog.Error("can't make a new module builder", "mod name", arg.algo.ModuleName, "algo", arg.algo.AlgorithmName, "err", err)
		return &scheduler.ResponseCorrect{}, errors.New("can't make a new module builder"), ""
	}
	initializationFunctions := strings.Split(arg.algo.Init, ",")
	modConfig := wazero.NewModuleConfig().WithStartFunctions()
	if len(initializationFunctions) != 0 {
		modConfig = modConfig.WithStartFunctions(initializationFunctions[0])
	}
	mod, err := runtime.InstantiateWithConfig(ctx, arg.algo.WasmBytes, modConfig)
	if err != nil {
		slog.Error("can't instantiate", "algo", arg.algo.AlgorithmName, "modConfig", modConfig, "err", err)
		return &scheduler.ResponseCorrect{}, errors.New("can't instantiate"), ""
	}
	for _, initFunction := range initializationFunctions[1:] {
		slog.Info("Calling init function", "init", initFunction == "my-init", "mod", mod == nil, "fun", mod.ExportedFunction("my-init") == nil, "add-card", mod.ExportedFunction("add-card") == nil)
		_, err = mod.ExportedFunction(initFunction).Call(ctx)
		if err != nil {
			slog.Error("Unable to call init function", "initFunction", initFunction, "algo", arg.algo.AlgorithmName, "err", err)
			return &scheduler.ResponseCorrect{}, errors.New("can't call init function"), ""
		}
	}
	getCard := mod.ExportedFunction("get-cards")
	malloc := mod.ExportedFunction(arg.algo.Alloc)
	if getCard == nil {
		slog.Error("get-card does not exists")
		return &scheduler.ResponseCorrect{}, errors.New("get-card does not exists"), ""
	}
	if malloc == nil {
		slog.Error("malloc is null", "func", arg.algo.Alloc)
		return &scheduler.ResponseCorrect{}, errors.New("malloc is null"), ""
	}
	free := mod.ExportedFunction(arg.algo.Dealloc)
	if free == nil {
		slog.Error("free is null", "func", arg.algo.Dealloc)
		return &scheduler.ResponseCorrect{}, errors.New("free is null"), ""
	}
	cards := map[uint32]*scheduler.Card{}

	slog.Info("cards len", "arg.allGrades", arg.allGrades)
	for _, grade := range arg.allGrades {
		realId := uint32(grade.QuestionID)
		slog.Info("grade", "realId", realId)
		v, exists := cards[realId]
		if !exists {
			cards[realId] = &scheduler.Card{
				Id: uint32(grade.QuestionID),
				CardRecords: []*scheduler.CardRecord{
					{
						When:    &timestamppb.Timestamp{Seconds: grade.Timestamp},
						Correct: grade.Result == 1,
					},
				},
				Tags: arg.tagsByQuestion[uint32(grade.QuestionID)],
			}
		} else {
			v.CardRecords = append(v.CardRecords, &scheduler.CardRecord{
				When:    &timestamppb.Timestamp{Seconds: grade.Timestamp},
				Correct: grade.Result == 1,
			})
		}
	}

	for id, tags := range arg.tagsByQuestion {
		if _, found := cards[id]; !found {
			cards[id] = &scheduler.Card{
				Id:          id,
				CardRecords: []*scheduler.CardRecord{},
				Tags:        tags,
			}
		}
	}

	var passCards []*scheduler.Card
	for _, card := range cards {
		passCards = append(passCards, card)
	}

	slog.Info("passCards", "p", passCards)
	slog.Info("passCards", "arg.tagsToAsk", arg.tagsToAsk)

	cards[0] = &scheduler.Card{}
	sched := scheduler.ToScheduler{
		CustomParams: `{}`,
		Cards:        passCards,
		TagsToQuery:  arg.tagsToAsk,
	}
	res, err := proto.Marshal(&sched)
	if err != nil {
		slog.Error("marshal err", "err", err)
		return &scheduler.ResponseCorrect{}, errors.New("marshalling error"), ""
	}
	slog.Info("marshalled result", "result", string(res))
	retPtrToWasm, err := malloc.Call(ctx, uint64(len(res)))
	if err != nil {
		slog.Error("malloc call err", "err", err)
		return &scheduler.ResponseCorrect{}, errors.New("malloc call err"), ""
	}
	ptrToWasm := retPtrToWasm[0]
	if valid := mod.Memory().Write(uint32(ptrToWasm), res); !valid {
		slog.Error("memory write err")
		return &scheduler.ResponseCorrect{}, errors.New("memory write err"), ""
	}
	ptrOutArr, err := getCard.Call(ctx, ptrToWasm, uint64(len(res)), uint64(time.Now().UTC().Unix()))
	if err != nil {
		slog.Error("getCard call err", "err", err)
		return &scheduler.ResponseCorrect{}, errors.New("getCard call err"), ""
	}
	p, l := splitU64ToPLen(ptrOutArr[0])
	slog.Info("called splitU64", "p", p, "l", l)
	m := readBytes(l, mod.Memory(), p)
	var retVal scheduler.FromScheduler
	slog.Info("called readBytes", "m", m)
	err = proto.Unmarshal(m, &retVal)
	if err != nil {
		slog.Error("unmarshal err from get-card", "err", err)
		return &scheduler.ResponseCorrect{}, errors.New("unmarshal err from get-card"), ""
	}
	if msg := retVal.GetBadValue(); msg != nil {
		slog.Error("spacing algorithm returned error", "log", msg.ToLog, "user", msg.ToUser)
		return &scheduler.ResponseCorrect{}, nil, msg.ToUser
	} else if msg := retVal.GetGoodValue(); msg != nil {
		return msg, nil, ""
	} else {
		slog.Info("Invalid retval", "retval type", retVal.String())
		return &scheduler.ResponseCorrect{}, nil, ""
	}
}

func splitU64ToPLen(val uint64) (uint32, uint32) {
	high := uint32(val >> 32)
	low := uint32(val & 0xffffffff)
	return high, low
}

func readBytes(l uint32, mem api.Memory, p uint32) []byte {
	m, valid := mem.Read(p, l)
	if !valid {
		panic("not valid read")
	}
	return m
}

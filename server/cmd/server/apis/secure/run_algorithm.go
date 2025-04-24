package secure

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"github.com/tetratelabs/wazero"
	"github.com/tetratelabs/wazero/api"
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

type AlgoReturn struct {
	lenNewCards    int
	lenDueCards    int
	lenNonDueCards int
	nextCard       int
	typeOfNextCard string
	isShortTerm    bool
}

func runAlgorithm(ctx context.Context, algo AlgorithmStruct, allQuestions []sql_queries.GetAllQuestionsRow, allGrades []sql_queries.QuestionsLog) (AlgoReturn, error) {
	v := sha256.Sum256(algo.WasmBytes)
	slog.Info("sha256sum", "sha", hex.EncodeToString(v[:]))
	runtime := wazero.NewRuntime(ctx)

	defer runtime.Close(ctx) // This closes everything this Runtime created.
	_, err := runtime.NewHostModuleBuilder(algo.ModuleName).Instantiate(ctx)
	if err != nil {
		slog.Error("can't make a new module builder", "mod name", algo.ModuleName, "algo", algo.AlgorithmName, "err", err)
		return AlgoReturn{}, errors.Join(ErrCantGetNewModBuilder, err)
	}
	initializationFunctions := strings.Split(algo.Init, ",")
	modConfig := wazero.NewModuleConfig().WithStartFunctions()
	if len(initializationFunctions) != 0 {
		modConfig = modConfig.WithStartFunctions(initializationFunctions[0])
	}
	mod, err := runtime.InstantiateWithConfig(ctx, algo.WasmBytes, modConfig)
	if err != nil {
		slog.Error("can't instantiate", "algo", algo.AlgorithmName, "modConfig", modConfig, "err", err)
		return AlgoReturn{}, errors.Join(ErrCantInstantiate, err)
	}
	for _, initFunction := range initializationFunctions[1:] {
		slog.Info("Calling init function", "init", initFunction == "my-init", "mod", mod == nil, "fun", mod.ExportedFunction("my-init") == nil, "add-card", mod.ExportedFunction("add-card") == nil)
		_, err = mod.ExportedFunction(initFunction).Call(ctx)
		if err != nil {
			slog.Error("Unable to call init function", "initFunction", initFunction, "algo", algo.AlgorithmName, "err", err)
			return AlgoReturn{}, errors.Join(ErrCantInstantiate, err)
		}
	}
	addCard := mod.ExportedFunction("add-card")
	gradeCard := mod.ExportedFunction("grade-card")
	getCard := mod.ExportedFunction("get-cards")
	malloc := mod.ExportedFunction(algo.Alloc)
	if getCard == nil {
		panic("getCard is null")
	}
	if malloc == nil {
		slog.Error("malloc is null", "func", algo.Alloc)
		return AlgoReturn{}, errors.Join(ErrCantInstantiate, err)
	}
	free := mod.ExportedFunction(algo.Dealloc)

	for _, question := range allQuestions {
		slog.Info("adding question", "id", question.QuestionID, "time", time.Now().Unix())
		_, err = addCard.Call(ctx, uint64(question.QuestionID))
		if err != nil {
			slog.Error("Unable to call add-card function", "algo", algo.AlgorithmName, "err", err)
			return AlgoReturn{}, errors.Join(ErrCantGetAddCard, err)
		}
	}
	for _, grade := range allGrades {
		slog.Info("grading question", "quid", uint64(grade.QuestionID), "time", uint64(grade.Timestamp), "grade", uint64(grade.Result))
		_, err := gradeCard.Call(ctx, uint64(grade.QuestionID), uint64(grade.Timestamp), uint64(grade.Result))
		if err != nil {
			slog.Error("Unable to call grade-card function", "algo", algo.AlgorithmName, "err", err)
			return AlgoReturn{}, errors.Join(ErrCantGetGradeCard, err)
		}
	}
	addr, err := malloc.Call(ctx, 22)
	if err != nil {
		slog.Error("Unable to allocate", "algo", algo.AlgorithmName, "err", err)
		return AlgoReturn{}, errors.Join(ErrCantAllocate, err)
	}

	_, err = getCard.Call(ctx, addr[0], uint64(time.Now().Unix()))
	if err != nil {
		slog.Error("Unable to get next card", "algo", algo.AlgorithmName, "time", time.Now().Unix(), "err", err)
		return AlgoReturn{}, errors.Join(ErrCantCallNextCard, err)
	}
	p := addr[0]
	lenDueCards := binary.LittleEndian.Uint32(readBytes(4, mod.Memory(), p))
	lenNonDueCards := binary.LittleEndian.Uint32(readBytes(4, mod.Memory(), p+4))
	lenNewCards := binary.LittleEndian.Uint32(readBytes(4, mod.Memory(), p+8))
	nextType := readBytes(1, mod.Memory(), p+12)[0]
	isShortTerm := false
	if readBytes(1, mod.Memory(), p+13)[0] == 1 {
		isShortTerm = true
	}
	nextCard := binary.LittleEndian.Uint64(readBytes(8, mod.Memory(), p+14))

	slog.Info("Next Card", "id", nextCard)
	_, err = free.Call(ctx, addr[0])
	if err != nil {
		slog.Error("Cannot call free", "algo", algo.AlgorithmName, "err", err)
		return AlgoReturn{}, ErrCantReadBytes
	}
	typeOfCardString := ""
	switch nextType {
	case 1:
		typeOfCardString = "due-card"
	case 2:
		typeOfCardString = "new-card"
	case 3:
		typeOfCardString = "non-due-card"
	}

	return AlgoReturn{
		lenNewCards:    int(lenNewCards),
		lenDueCards:    int(lenDueCards),
		lenNonDueCards: int(lenNonDueCards),
		nextCard:       int(nextCard),
		isShortTerm:    bool(isShortTerm),
		typeOfNextCard: typeOfCardString,
	}, nil
}

func readBytes(l int, mem api.Memory, p uint64) []byte {
	m, valid := mem.Read(uint32(p), uint32(l))
	if !valid {
		panic("not valid read")
	}
	return m
}

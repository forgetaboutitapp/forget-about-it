package secure

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"github.com/tetratelabs/wazero"
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
	typeOfNextCard int
}

func runAlgorithm(ctx context.Context, algo AlgorithmStruct, allQuestions []sql_queries.GetAllQuestionsRow, allGrades []sql_queries.QuestionsLog) (AlgoReturn, error) {

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
		modConfig = modConfig.WithStartFunctions(initializationFunctions...)
	}
	mod, err := runtime.InstantiateWithConfig(ctx, algo.WasmBytes, modConfig)
	if err != nil {
		slog.Error("can't instantiate", "algo", algo.AlgorithmName, "modConfig", modConfig, "err", err)
		return AlgoReturn{}, errors.Join(ErrCantInstantiate, err)
	}
	addCard := mod.ExportedFunction("add-card")
	gradeCard := mod.ExportedFunction("grade-card")
	getCard := mod.ExportedFunction("get-cards")
	malloc := mod.ExportedFunction("alloc")
	free := mod.ExportedFunction("dealloc")

	for _, question := range allQuestions {
		_, err = addCard.Call(ctx, uint64(question.QuestionID))
		if err != nil {
			slog.Error("Unable to call add-card function", "algo", algo.AlgorithmName, "err", err)
			return AlgoReturn{}, errors.Join(ErrCantGetAddCard, err)
		}
	}
	for _, grade := range allGrades {
		_, err := gradeCard.Call(ctx, uint64(grade.QuestionID), uint64(grade.Timestamp), uint64(grade.Result))
		if err != nil {
			slog.Error("Unable to call grade-card function", "algo", algo.AlgorithmName, "err", err)
			return AlgoReturn{}, errors.Join(ErrCantGetGradeCard, err)
		}
	}
	addr, err := malloc.Call(ctx, 25)
	if err != nil {
		slog.Error("Unable to allocate", "algo", algo.AlgorithmName, "err", err)
		return AlgoReturn{}, errors.Join(ErrCantAllocate, err)
	}
	_, err = getCard.Call(ctx, addr[0], uint64(time.Now().Unix()))
	if err != nil {
		slog.Error("Unable to get next card", "algo", algo.AlgorithmName, "err", err)
		return AlgoReturn{}, errors.Join(ErrCantCallNextCard, err)
	}

	lenDueCards, inRange := mod.Memory().ReadUint32Le(uint32(addr[0]))
	if !inRange {
		slog.Error("Cannot read 4 bytes in getting amount of due cards", "algo", algo.AlgorithmName)
		return AlgoReturn{}, ErrCantReadBytes
	}

	lenNonDueCards, inRange := mod.Memory().ReadUint32Le(uint32(addr[0] + 4))
	if !inRange {
		slog.Error("Cannot read 4 bytes in getting amount of non due cards", "algo", algo.AlgorithmName)
		return AlgoReturn{}, ErrCantReadBytes
	}

	lenNewCards, inRange := mod.Memory().ReadUint32Le(uint32(addr[0] + 8))
	if !inRange {
		slog.Error("Cannot read 4 bytes in getting amount of new cards", "algo", algo.AlgorithmName)
		return AlgoReturn{}, ErrCantReadBytes
	}
	nextCard, inRange := mod.Memory().ReadUint64Le(uint32(addr[0] + 16))
	if !inRange {
		slog.Error("Cannot read 8 bytes in getting due card", "algo", algo.AlgorithmName)
		return AlgoReturn{}, ErrCantReadBytes
	}

	nextType, inRange := mod.Memory().ReadByte(uint32(addr[0] + 25))
	if !inRange {
		slog.Error("Cannot read 1 byte in getting next card type", "algo", algo.AlgorithmName)
		return AlgoReturn{}, ErrCantReadBytes
	}
	_, err = free.Call(ctx, addr[0])
	if err != nil {
		slog.Error("Cannot call free", "algo", algo.AlgorithmName, "err", err)
		return AlgoReturn{}, ErrCantReadBytes
	}
	return AlgoReturn{
		lenNewCards:    int(lenNewCards),
		lenDueCards:    int(lenDueCards),
		lenNonDueCards: int(lenNonDueCards),
		nextCard:       int(nextCard),
		typeOfNextCard: int(nextType),
	}, nil
}

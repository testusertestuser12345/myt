package commands

import (
	"fmt"
	"strconv"

	structs "github.com/MythicMeta/MythicContainer/agent_structs"
	"github.com/MythicMeta/MythicContainer/logging"

	"github.com/Ne0nd0g/merlin-message/jobs"
)

func tcpdumpSimple() structs.Command {
	attr := structs.CommandAttribute{
		SupportedOS: []string{structs.SUPPORTED_OS_LINUX, structs.SUPPORTED_OS_MACOS, structs.SUPPORTED_OS_WINDOWS},
	}

	// UIのフォームに表示するパラメータ（キャプチャするパケット数）
	packetCount := structs.CommandParameter{
		Name:             "count",
		ModalDisplayName: "Packet Count",
		CLIName:          "count",
		ParameterType:    structs.COMMAND_PARAMETER_TYPE_NUMBER, // 数値型
		Description:      "Number of packets to capture before returning results",
		ParameterGroupInformation: []structs.ParameterGroupInfo{
			{ParameterIsRequired: true, GroupName: "Default", UIModalPosition: 0},
		},
	}

	command := structs.Command{
		Name:                           "tcpdump_simple",
		NeedsAdminPermissions:          true, // パケットキャプチャやRAWソケットには管理者権限が必要
		HelpString:                     "tcpdump_simple 5",
		Description:                    "Capture a specified number of packets using tcpdump or native fallback",
		Version:                        0,
		Author:                         "@you",
		MitreAttackMappings:            []string{"T1040"},
		ScriptOnlyCommand:              false,
		CommandAttributes:              attr,
		CommandParameters:              []structs.CommandParameter{packetCount},
		TaskFunctionCreateTasking:      tcpdumpSimpleCreateTask,
		TaskFunctionParseArgString:     taskFunctionParseArgString,
		TaskFunctionParseArgDictionary: taskFunctionParseArgDictionary,
	}
	return command
}

// 実行時に頭脳側で走る処理：UIから値を受け取り、エージェント（体）へ渡すジョブを組み立てる
func tcpdumpSimpleCreateTask(task *structs.PTTaskMessageAllData) (resp structs.PTTaskCreateTaskingMessageResponse) {
	pkg := "mythic/container/commands/tcpdump_simple/tcpdumpSimpleCreateTask()"
	resp.TaskID = task.Task.ID

	count, err := task.Args.GetNumberArg("count") // UIから数値を取得 (float64)
	if err != nil {
		err = fmt.Errorf("%s: there was an error getting the 'count' argument: %s", pkg, err)
		resp.Error = err.Error()
		resp.Success = false
		logging.LogError(err, "returning with error")
		return
	}

	countStr := strconv.Itoa(int(count)) // 体側へ渡すために文字列に変換

	job := jobs.Command{
		Command: task.Task.CommandName, // = "tcpdump_simple"
		Args:    []string{countStr},    // = ["5"] → 体側の cmd.Args[0]
	}

	mythicJob, err := ConvertMerlinJobToMythicTask(job, jobs.NATIVE)
	if err != nil {
		err = fmt.Errorf("%s: there was an error converting the Merlin job to a Mythic task: %s", pkg, err)
		resp.Error = err.Error()
		resp.Success = false
		logging.LogError(err, "returning with error")
		return
	}

	task.Args.SetManualArgs(mythicJob)

	resp.DisplayParams = &countStr // タスク一覧に引数を表示
	resp.Success = true
	return
}

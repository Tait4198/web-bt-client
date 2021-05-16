package db

import (
	"context"
	"crawshaw.io/sqlite"
	"crawshaw.io/sqlite/sqlitex"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type Task struct {
	InfoHash           string   `json:"info_hash"`
	TorrentName        string   `json:"torrent_name"`
	Complete           bool     `json:"complete"`
	MetaInfo           bool     `json:"meta_info"`
	Pause              bool     `json:"pause"`
	Download           bool     `json:"download"`
	DownloadPath       string   `json:"download_path"`
	DownloadFiles      []string `json:"download_files"`
	FileLength         int64    `json:"file_length"`
	CompleteFileLength int64    `json:"complete_file_length"`
	CreateTime         int64    `json:"create_time"`
	CompleteTime       int64    `json:"complete_time"`
	CreateTorrentInfo  string   `json:"create_torrent_info"`
}

func SelectTaskCount(infoHash string) int {
	conn := GetPool().Get(context.TODO())
	defer GetPool().Put(conn)
	stmt := conn.Prep("select count(*) from tasks where info_hash = $hash")
	stmt.SetText("$hash", infoHash)
	if count, err := sqlitex.ResultInt(stmt); err == nil {
		return count
	} else {
		log.Println(fmt.Errorf("获取 Task %s 数量失败 %w \n", infoHash, err))
	}
	return -1
}

func SelectActiveTaskList() ([]Task, error) {
	return SelectTaskListBase("select * from tasks where pause = 0")
}

func SelectTaskList() ([]Task, error) {
	return SelectTaskListBase("select * from tasks order by create_time desc")
}

func SelectTask(infoHash string) (Task, error) {
	if tasks, err := SelectTaskListBase(fmt.Sprintf("select * from tasks where info_hash = '%s'", infoHash)); err == nil {
		if len(tasks) != 1 {
			return Task{}, fmt.Errorf("SelectTask %s 不存在", infoHash)
		} else {
			return tasks[0], nil
		}
	} else {
		return Task{}, err
	}
}

func SelectTaskListBase(sql string) ([]Task, error) {
	conn := GetPool().Get(context.TODO())
	defer GetPool().Put(conn)
	stmt := conn.Prep(sql)

	defer func() {
		if err := stmt.Reset(); err != nil {
			log.Println(fmt.Errorf("GetTaskList stmt.Reset() 失败 %w \n", err))
		}
	}()

	var tasks []Task
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, fmt.Errorf("查询任务失败 %w", err)
		} else if !hasRow {
			break
		}
		if task, err := stmtConvertTask(stmt); err != nil {
			return nil, err
		} else {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}

func InsertTask(task Task) error {
	if task.DownloadFiles == nil {
		task.DownloadFiles = []string{}
	}
	downloadFiles, err := json.Marshal(task.DownloadFiles)
	if err != nil {
		return fmt.Errorf("InsertTask Marshal 失败 %w", err)
	}
	return ExecSql("insert into tasks values (?,?,?,?,?,?,?,?,?,?,?,?,?)",
		task.InfoHash,
		task.TorrentName,
		boolToInt(task.Complete),
		boolToInt(task.MetaInfo),
		boolToInt(task.Pause),
		boolToInt(task.Download),
		task.DownloadPath,
		downloadFiles,
		task.FileLength,
		task.CompleteFileLength,
		task.CreateTime,
		task.CompleteTime,
		task.CreateTorrentInfo)
}

func UpdateTaskMetaInfo(infoHash, torrentName string, fileLength int64) error {
	// todo 增加WS推送
	return ExecSql("update tasks set torrent_name = ?,file_length = ?,meta_info = 1 where info_hash = ?",
		torrentName, fileLength, infoHash)
}

func UpdateTaskPause(pause bool, infoHash string) error {
	if pause {
		return ExecSql("update tasks set pause = ? where info_hash = ?", boolToInt(pause), infoHash)
	} else {
		return ExecSql("update tasks set pause = ?, complete = 0 where info_hash = ?", boolToInt(pause), infoHash)
	}
}

func UpdateTaskDownload(download bool, infoHash string) error {
	return ExecSql("update tasks set download = ? where info_hash = ?", boolToInt(download), infoHash)
}

func UpdateTaskComplete(complete bool, infoHash string) error {
	return ExecSql("update tasks set complete = ? where info_hash = ?", boolToInt(complete), infoHash)
}

func UpdateTaskDownloadFiles(files []string, infoHash string) error {
	downloadFiles, err := json.Marshal(files)
	if err != nil {
		return fmt.Errorf("UpdateTaskDownloadFiles Marshal 失败 %w", err)
	}
	return ExecSql("update tasks set download_files = ? where info_hash = ?", downloadFiles, infoHash)
}

func UpdateTaskCompleteFileLength(completeFileLength int64, infoHash string) error {
	return ExecSql("update tasks set complete_file_length = ? where info_hash = ?",
		completeFileLength, infoHash)
}

func TaskComplete(completeFileLength int64, infoHash string) error {
	return ExecSql("update tasks set complete = ?, complete=1, download = 0, pause = 1, complete_time = ? where info_hash = ?",
		completeFileLength, time.Now().UnixNano(), infoHash)
}

func stmtConvertTask(stmt *sqlite.Stmt) (Task, error) {
	task := Task{}
	task.InfoHash = stmt.GetText("info_hash")
	task.TorrentName = stmt.GetText("torrent_name")
	task.Complete = stmt.GetInt64("complete") == 1
	task.MetaInfo = stmt.GetInt64("meta_info") == 1
	task.Pause = stmt.GetInt64("pause") == 1
	task.Download = stmt.GetInt64("download") == 1
	task.DownloadPath = stmt.GetText("download_path")
	task.FileLength = stmt.GetInt64("file_length")
	task.CompleteFileLength = stmt.GetInt64("complete_file_length")
	task.CreateTime = stmt.GetInt64("create_time")
	task.CompleteTime = stmt.GetInt64("complete_time")
	task.CreateTorrentInfo = stmt.GetText("create_torrent_info")

	var downloadFiles []string
	err := json.Unmarshal([]byte(stmt.GetText("download_files")), &downloadFiles)
	if err != nil {
		return task, fmt.Errorf("GetTaskList Unmarshal 失败 %w", err)
	}
	task.DownloadFiles = downloadFiles
	return task, nil
}

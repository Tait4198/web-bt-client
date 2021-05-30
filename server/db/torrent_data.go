package db

import (
	"context"
	"crawshaw.io/sqlite/sqlitex"
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"log"
	"strings"
)

func SelectMetaInfo(infoHash string) (*metainfo.MetaInfo, error) {
	if mis, err := SelectMateInfoList([]string{infoHash}); err == nil && len(mis) == 1 {
		return mis[0], nil
	} else if len(mis) == 0 {
		return nil, fmt.Errorf("SelectMetaInfo %s 不存在", infoHash)
	} else {
		return nil, err
	}
}

func SelectMateInfoList(infoHash []string) ([]*metainfo.MetaInfo, error) {
	conn := GetPool().Get(context.TODO())
	defer GetPool().Put(conn)

	sql := fmt.Sprintf("select * from torrent_data where info_hash in (%s)",
		fmt.Sprintf("'%s'", strings.Join(infoHash, "','")))
	stmt := conn.Prep(sql)
	defer func() {
		if err := stmt.Reset(); err != nil {
			log.Println(fmt.Errorf("GetMetaInfo stmt.Reset() 失败 %w \n", err))
		}
	}()
	var mis []*metainfo.MetaInfo
	for {
		if hasRow, err := stmt.Step(); err != nil {
			return nil, err
		} else if !hasRow {
			break
		}
		r := stmt.GetReader("torrent_blob")
		if mi, err := metainfo.Load(r); err == nil {
			mis = append(mis, mi)
		} else {
			return nil, fmt.Errorf("bytes.Reader 转换 MetaInfo 失败 %w \n", err)
		}
	}
	return mis, nil
}

func SelectTorrentDataCount(infoHash string) int {
	conn := GetPool().Get(context.TODO())
	defer GetPool().Put(conn)
	stmt := conn.Prep("select count(*) from torrent_data where info_hash = $hash")
	stmt.SetText("$hash", infoHash)
	if count, err := sqlitex.ResultInt(stmt); err == nil {
		return count
	} else {
		log.Println(fmt.Errorf("获取 Hash %s 数量失败 %w \n", infoHash, err))
	}
	return -1
}

func InsertTorrentData(hash string, blob []byte) error {
	return ExecSql("insert into torrent_data values (?,?);", hash, blob)
}

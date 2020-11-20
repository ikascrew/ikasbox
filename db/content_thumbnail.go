package db

import "golang.org/x/xerrors"

//Argen で[]byteが使えないぽいので個別に書き込む
type ContentThumbnail struct {
	ID   int
	Seq  int
	Data []byte
}

func (t *ContentThumbnail) Insert() error {
	sq := "insert into content_thumbnails(id,seq,data) values (?,?,?)"
	_, err := db.Exec(sq, t.ID, t.Seq, t.Data)
	if err != nil {
		return xerrors.Errorf("content thumbnail insert: %w", err)
	}
	return nil
}

func SelectContentThumbnails(ID int) ([]*ContentThumbnail, error) {

	sq := "select id,seq,data from content_thumbnails order by seq"

	rows, err := db.Query(sq)
	if err != nil {
		return nil, xerrors.Errorf("query error : %w", err)
	}
	defer rows.Close()

	data := make([]*ContentThumbnail, 0)

	for rows.Next() {
		row := ContentThumbnail{}
		row.Data = make([]byte, 0)
		err = rows.Scan(&row.ID, &row.Seq, &row.Data)
		if err != nil {
			return nil, xerrors.Errorf("rows scan error : %w", err)
		}
		data = append(data, &row)
	}

	return data, nil
}

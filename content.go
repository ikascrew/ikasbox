package ikasbox

import (
	"fmt"
	"strconv"

	"github.com/ikascrew/ikasbox/config"
	"github.com/ikascrew/ikasbox/contentimport"
	"github.com/ikascrew/ikasbox/db"
	"github.com/ikascrew/plugin/video"

	"golang.org/x/xerrors"
)

func setContent() error {

	err := db.Open()
	if err != nil {
		return xerrors.Errorf("Database Open : %w", err)
	}

	conf := config.Get()

	switch conf.Function {
	case "register":
		err = registerGeneratedContent(conf.Arguments)
	default:
		err = fmt.Errorf("not found function")
	}

	if err != nil {
		return xerrors.Errorf("function error[%s]: %w", conf.Function, err)
	}

	return nil
}

// registerGeneratedContent は実体ファイルを持たない生成型コンテンツを登録する。
//
//	ikasbox content register <group-id> <name> <type> [params-json]
//	例: content register 1 "NewYear" cd "{\"target\":\"2027-01-01T00:00:00+09:00\",\"text\":\"HappyNewYear\"}"
//
// グループ ID は「group list」で確認できる
func registerGeneratedContent(args []string) error {

	if len(args) < 3 {
		return fmt.Errorf("usage: content register <group-id> <name> <type> [params-json]\n  types: %v", video.Types())
	}

	id, err := strconv.Atoi(args[0])
	if err != nil {
		return xerrors.Errorf("group id[%s]: %w", args[0], err)
	}

	name := args[1]
	typ := args[2]
	params := ""
	if len(args) > 3 {
		params = args[3]
	}

	c, err := contentimport.RegisterGenerated(id, name, typ, params)
	if err != nil {
		return xerrors.Errorf("register generated: %w", err)
	}

	fmt.Printf("New Content:%s[%d] type=%s params=%s\n", c.Name, c.ID, c.Type, c.Params)
	return nil
}

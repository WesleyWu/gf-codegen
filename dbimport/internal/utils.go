package internal

import (
	"context"
	_ "embed"
	"github.com/WesleyWu/gf-codegen/model"
	"github.com/WesleyWu/gf-codegen/util"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"os"
	"path"
)

//go:embed template/yaml/table_def.template
var yamlTemplate string

func ParseDblink(dblink string) error {
	match, _ := gregex.MatchString(`([a-z]+):(.+)`, dblink)
	if len(match) == 3 {
		gdb.AddDefaultConfigNode(gdb.ConfigNode{
			Type: gstr.Trim(match[1]),
			Link: gstr.Trim(match[2]),
		})
		return nil
	}
	return gerror.Newf("不正确的dblink格式：%s", dblink)
}

func SaveTableDef(ctx context.Context, table *model.TableDef, yamlOutputPath string) error {
	curDir, err := os.Getwd()
	if err != nil {
		return gerror.New("获取本地路径失败")
	}
	yamlFile := path.Join(curDir, yamlOutputPath, table.Name+".yaml")

	view := util.TemplateEngine()
	tplData := g.Map{"apiVersion": "v1", "table": table}
	var tplOut string
	if tplOut, err = view.ParseContent(ctx, yamlTemplate, tplData); err != nil {
		return err
	}
	tplOut, err = util.TrimBreak(tplOut)
	if err != nil {
		return err
	}
	err = util.WriteFile(yamlFile, tplOut, true)
	if err != nil {
		return err
	}
	return nil
}

func SplitComma(str string) []string {
	var result []string
	if g.IsEmpty(gstr.Trim(str)) {
		return result
	}
	temp := gstr.Split(str, ",")
	for _, one := range temp {
		oneTrimmed := gstr.Trim(one)
		if !g.IsEmpty(oneTrimmed) {
			result = append(result, oneTrimmed)
		}
	}
	return result
}

// GetDbDriver 获取数据库驱动类型
func GetDbDriver() string {
	config := g.DB(gdb.DefaultGroupName).GetConfig()
	return gstr.ToLower(config.Type)
}

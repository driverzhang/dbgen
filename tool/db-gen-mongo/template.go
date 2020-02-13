package db_gen_mongo

import (
	"bytes"
	"text/template"
)

// VersionOptions include version
type TableOptions struct {
	N       string
	LowName string
	Name    string
}

func (t *TableOptions) setName(name string) {
	t.LowName = name
}

func (t *TableOptions) getMongoCrudTemplate() (rsp string, err error) {
	var doc bytes.Buffer
	tmpl, _ := template.New("TableOptions").Parse(crudTemplate)
	err = tmpl.Execute(&doc, t)
	rsp = doc.String()
	return
}

var crudTemplate = `func ({{.N}} *{{.Name}}) TableName() (rsp string) {
	return "{{.LowName}}"
}

func ({{.N}} *{{.Name}}) CreateManyIndex() (err error) {
	keys := map[string]interface{}{}
	err = mongo.Collection({{.N}}.TableName()).CreateManyIndex(keys)
	if err != nil {
		return
	}

	return
}

func ({{.N}} *{{.Name}}) CreateIndex() (err error) {
	keys := map[string]interface{}{}
	err = mongo.Collection({{.N}}.TableName()).CreateUniqueIndex(keys)
	if err != nil {
		return
	}

	return
}

func ({{.N}} *{{.Name}}) GetList(query, sort bson.M, from, size int) (rsp []{{.Name}}, err error) {
	collection := mongo.Collection({{.N}}.TableName())
	if sort == nil {
		sort = bson.M{"created_at": 1}
	}
	err = collection.Where(query).Sort(sort).Skip(int64(from)).Limit(int64(size)).FindMany2(&rsp)
	if err != nil {
		err = errors.New("获取{{.Name}}列表出错！服务器繁忙 " + err.Error())
		return
	}

	return
}

func ({{.N}} *{{.Name}}) GetOne(conditions bson.M) (err error) {
	err = mongo.Collection({{.N}}.TableName()).Where(conditions).FindOne({{.N}})
	if err != nil {
		err = errors.New("获取{{.Name}}详情出错！" + err.Error())
		return
	}
	return
}

func ({{.N}} *{{.Name}}) UpdateOne(data bson.M) (updateId interface{}, err error) {
	if {{.N}}.Id.IsZero() {
		err = errors.New("invalid _id")
		return
	}

	collection := mongo.Collection({{.N}}.TableName())
	where := bson.M{"_id": {{.N}}.Id}
	_, err = collection.Where(where).UpdateOne2(data)
	if err != nil {
		err = errors.New("更新{{.Name}}详情错误！服务器繁忙 " + err.Error())
		return
	}
	updateId = {{.N}}.Id
	return
}

func ({{.N}} *{{.Name}}) InsertOne() (insertId interface{}, err error) {
	collection := mongo.Collection({{.N}}.TableName())
	insertR, err := collection.InsertOne2({{.N}})
	if err != nil {
		err = errors.New("生成{{.Name}}数据错误！服务器繁忙 " + err.Error())
		return
	}
	insertId = insertR.InsertedID
	return
}

func ({{.N}} *{{.Name}}) DeleteMany(ids []string) (count int64, err error) {
	objIds := make([]primitive.ObjectID, len(ids))
	for i, v := range ids {
		objId, _ := primitive.ObjectIDFromHex(v)
		objIds[i] = objId
	}

	collection := mongo.Collection({{.N}}.TableName())
	where := bson.M{"_id": bson.M{"$in": objIds}}
	count, err = collection.Where(where).Delete2()
	if err != nil {
		err = errors.New("删除{{.Name}}信息错误！服务器繁忙 " + err.Error())
		return
	}

	return
}

func ({{.N}} *{{.Name}}) Count(statement bson.M) (count int64, err error) {
	count, err = mongo.Collection({{.N}}.TableName()).Where(statement).Count2()
	if err != nil {
		err = errors.New("获取{{.Name}}数据条数出错！服务器繁忙 " + err.Error())
		return
	}
	return
}`

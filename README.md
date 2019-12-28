# go-mongo-crud-template
this is a mongo CRUD template tool for golang

# how to use


copy some struct name to your clipboard

> 首先保证你已经复制了一个有效的struct名称到你的剪贴板中

    dbgen mongo

or 

    dbgen mo
    
if you can see :

    gen mongo db-gen-mongo template struct from clipboard success
   
so now you can paste some codes to your IDE


# demo:

    type Back struct {}
  
so i can copy this "Back"

    dbgen mongo
    
my codes :

```go
func (b *Back) TableName() (rsp string) {
	return "back"
}

func (b *Back) CreateManyIndex() (err error) {
	keys := map[string]interface{}{}
	err = mongo.Collection(b.TableName()).CreateManyIndex(keys)
	if err != nil {
		return
	}

	return
}

func (b *Back) CreateIndex() (err error) {
	keys := map[string]interface{}{}
	err = mongo.Collection(b.TableName()).CreateUniqueIndex(keys)
	if err != nil {
		return
	}

	return
}

func (b *Back) GetList(query, sort bson.M, from, size int) (rsp []Back, err error) {
	collection := mongo.Collection(b.TableName())
	if sort == nil {
		sort = bson.M{"created_at": 1}
	}
	err = collection.Where(query).Sort(sort).Skip(int64(from)).Limit(int64(size)).FindMany2(&rsp)
	if err != nil {
		err = errors.New("获取Back列表出错！服务器繁忙 " + err.Error())
		return
	}

	return
}

func (b *Back) GetOne(conditions bson.M) (err error) {
	err = mongo.Collection(b.TableName()).Where(conditions).FindOne(f)
	if err != nil {
		err = errors.New("获取Back详情出错！" + err.Error())
		return
	}
	return
}

func (b *Back) UpdateOne(data bson.M) (updateId interface{}, err error) {
	if b.Id.IsZero() {
		err = errors.New("invalid _id")
		return
	}

	collection := mongo.Collection(b.TableName())
	where := bson.M{"_id": b.Id}
	_, err = collection.Where(where).UpdateOne2(data)
	if err != nil {
		err = errors.New("更新Back详情错误！服务器繁忙 " + err.Error())
		return
	}
	updateId = b.Id
	return
}

func (b *Back) InsertOne() (insertId interface{}, err error) {
	collection := mongo.Collection(b.TableName())
	insertR, err := collection.InsertOne2(b)
	if err != nil {
		err = errors.New("生成Back数据错误！服务器繁忙 " + err.Error())
		return
	}
	insertId = insertR.InsertedID
	return
}

func (b *Back) DeleteMany(ids []string) (count int64, err error) {
	objIds := make([]primitive.ObjectID, len(ids))
	for i, v := range ids {
		objId, _ := primitive.ObjectIDFromHex(v)
		objIds[i] = objId
	}

	collection := mongo.Collection(b.TableName())
	where := bson.M{"_id": bson.M{"$in": objIds}}
	count, err = collection.Where(where).Delete2()
	if err != nil {
		err = errors.New("删除Back信息错误！服务器繁忙 " + err.Error())
		return
	}

	return
}

func (b *Back) Count(statement bson.M) (count int64, err error) {
	count, err = mongo.Collection(b.TableName()).Where(statement).Count2()
	if err != nil {
		err = errors.New("获取Back数据条数出错！服务器繁忙 " + err.Error())
		return
	}
	return
}
```
  
 
 
 
 # todoList:
 
 - mysql db gen
 
package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/driverzhang/dbgen/pkg/config"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/driverzhang/dbgen/pkg/common/log"
)

var (
	client *mongo.Client
)

//创建mongo 索引
type MongoIndex interface {
	CreateIndex() error //创建mongo索引
}

type MgoUniqueIndex interface {
	CreateManyIndex() error
}

// 创建多字段index索引
func ManyIndexInit(values ...MgoUniqueIndex) {
	for _, v := range values {
		err := v.CreateManyIndex()
		if err != nil {
			log.Println(err)
		}
	}
}

// 创建单个唯一索引
func IndexInit(values ...MongoIndex) {
	for _, v := range values {
		err := v.CreateIndex()
		if err != nil {
			log.Println(err)
		}
	}
}

type collection struct {
	Database *mongo.Database
	Table    *mongo.Collection
	filter   bson.M
	limit    int64
	skip     int64
	sort     bson.M
	fields   bson.M
}

// 启动 mongo
func Start() {
	var err error
	mongoOptions := options.Client()
	mongoOptions.SetMaxConnIdleTime(time.Duration(config.Config.Mongo.MaxConnIdleTime) * time.Second)
	mongoOptions.SetMaxPoolSize(uint64(config.Config.Mongo.MaxPoolSize))
	if config.Config.Mongo.Username != "" && config.Config.Mongo.Password != "" {
		mongoOptions.SetAuth(options.Credential{Username: config.Config.Mongo.Username, Password: config.Config.Mongo.Password})
	}

	client, err = mongo.NewClient(mongoOptions.ApplyURI(config.Config.Mongo.Url))
	if err != nil {
		log.Fatalln(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

// 得到一个mongo操作对象
func Collection(table string) *collection {
	database := client.Database(config.Config.Mongo.Database)
	return &collection{
		Database: database,
		Table:    database.Collection(table),
		filter:   make(bson.M),
	}
}

func GetClient() *mongo.Client {
	return client
}

// 多字段创建index索引
func (c *collection) CreateManyIndex(keys map[string]interface{}) (err error) {
	ctx := context.Background()
	indexView := c.Table.Indexes()
	indexModels := make([]mongo.IndexModel, len(keys))
	j := 0
	for i, v := range keys {
		key := map[string]interface{}{i: v}
		indexModels[j] = mongo.IndexModel{
			Keys: key,
		}
		j++
	}

	res, err := indexView.CreateMany(ctx, indexModels)
	if err != nil {
		return
	}
	fmt.Printf("创建多字段索引：%+v\n", res)
	return
}

// 单字段创建唯一索引
func (collection *collection) CreateUniqueIndex(keys map[string]interface{}) error {
	ctx := context.Background()
	unique := true
	indexView := collection.Table.Indexes()
	option := options.Index()
	option.Unique = &unique
	indexModel := mongo.IndexModel{Keys: keys, Options: option}
	res, err := indexView.CreateOne(ctx, indexModel)
	if err != nil {
		return err
	}
	fmt.Printf("创建唯一索引：%+v\n", res)
	return nil
}

// 条件查询, bson.M{"field": "value"}
func (collection *collection) Where(m bson.M) *collection {
	collection.filter = m
	return collection
}

// 限制条数
func (collection *collection) Limit(n int64) *collection {
	collection.limit = n
	return collection
}

// 跳过条数
func (collection *collection) Skip(n int64) *collection {
	collection.skip = n
	return collection
}

// 排序 bson.M{"created_at":-1}
func (collection *collection) Sort(sorts bson.M) *collection {
	collection.sort = sorts
	return collection
}

// 指定查询字段
func (collection *collection) Fields(fields bson.M) *collection {
	collection.fields = fields
	return collection
}

func (collection *collection) InsertOne2(document interface{}) (result *mongo.InsertOneResult, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err = collection.Table.InsertOne(ctx, BeforeCreate(document))
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// 写入单条数据
// Deprecated: 将被删除，请调用 InsertOne2
func (collection *collection) InsertOne(document interface{}) *mongo.InsertOneResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.InsertOne(ctx, BeforeCreate(document))
	if err != nil {
		log.Println(err)
	}
	return result
}

func (collection *collection) InsertMany2(documents interface{}) (rsp *mongo.InsertManyResult, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var data []interface{}
	data = BeforeCreate(documents).([]interface{})
	rsp, err = collection.Table.InsertMany(ctx, data)
	if err != nil {
		log.Println(err)
		return
	}

	return
}

// 写入多条数据
// Deprecated: 将被删除，请调用 InsertMany2
func (collection *collection) InsertMany(documents interface{}) *mongo.InsertManyResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var data []interface{}
	data = BeforeCreate(documents).([]interface{})
	result, err := collection.Table.InsertMany(ctx, data)
	if err != nil {
		log.Println(err)
	}
	return result
}

// Deprecated: 将被删除，请调用 Aggregate2
func (collection *collection) Aggregate(pipeline interface{}, result interface{}) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := collection.Table.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err)
	}
	cursor.All(ctx, result)
}

func (collection *collection) Aggregate2(pipeline interface{}, result interface{}) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	cursor, err := collection.Table.Aggregate(ctx, pipeline)
	if err != nil {
		log.Println(err)
		return
	}
	err = cursor.All(ctx, result)
	if err != nil {
		return
	}
	return
}

// 存在更新,不存在写入, documents 里边的文档需要有 _id 的存在
func (collection *collection) UpdateOrInsert(documents interface{}) *mongo.UpdateResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	var upsert = true
	result, err := collection.Table.UpdateMany(ctx, collection.filter, documents, &options.UpdateOptions{Upsert: &upsert})
	if err != nil {
		log.Println(err)
	}
	return result
}

func (collection *collection) UpdateOne2(document interface{}) (result *mongo.UpdateResult, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err = collection.Table.UpdateOne(ctx, collection.filter, bson.M{"$set": BeforeUpdate(document)})
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// Deprecated: 将被删除，请调用 UpdateOne2
func (collection *collection) UpdateOne(document interface{}, opt ...*options.UpdateOptions) *mongo.UpdateResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.UpdateOne(ctx, collection.filter, bson.M{"$set": BeforeUpdate(document)}, opt...)
	if err != nil {
		log.Println(err)
	}
	return result
}

//原生update
func (collection *collection) UpdateOneRaw(document interface{}, opt ...*options.UpdateOptions) *mongo.UpdateResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.UpdateOne(ctx, collection.filter, document, opt...)
	if err != nil {
		log.Println(err)
	}
	return result
}

//
func (collection *collection) UpdateMany(document interface{}) *mongo.UpdateResult {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.UpdateMany(ctx, collection.filter, bson.M{"$set": BeforeUpdate(document)})
	if err != nil {
		log.Println(err)
	}
	return result
}

// 查询一条数据
func (collection *collection) FindOne(document interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result := collection.Table.FindOne(ctx, collection.filter, &options.FindOneOptions{
		Skip:       &collection.skip,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	err := result.Decode(document)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 查询多条数据，将错误外抛出
func (collection *collection) FindMany2(documents interface{}) (err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.Find(ctx, collection.filter, &options.FindOptions{
		Skip:       &collection.skip,
		Limit:      &collection.limit,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer result.Close(ctx)

	val := reflect.ValueOf(documents)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		log.Println("result argument must be a slice address")
		err = errors.New("result argument must be a slice address")
		return
	}

	slice := reflect.MakeSlice(val.Elem().Type(), 0, 0)
	itemTyp := val.Elem().Type().Elem()
	for result.Next(ctx) {
		item := reflect.New(itemTyp)
		err := result.Decode(item.Interface())
		if err != nil {
			log.Println(err)
			err = errors.New("result argument must be a slice address")
			return err
		}

		slice = reflect.Append(slice, reflect.Indirect(item))
	}
	val.Elem().Set(slice)
	return
}

// 查询多条数据
// Deprecated: 将被删除，请调用FindMany2
func (collection *collection) FindMany(documents interface{}) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.Find(ctx, collection.filter, &options.FindOptions{
		Skip:       &collection.skip,
		Limit:      &collection.limit,
		Sort:       collection.sort,
		Projection: collection.fields,
	})
	if err != nil {
		log.Println(err)
	}
	defer result.Close(ctx)

	val := reflect.ValueOf(documents)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Slice {
		log.Println("result argument must be a slice address")
	}

	slice := reflect.MakeSlice(val.Elem().Type(), 0, 0)

	itemTyp := val.Elem().Type().Elem()
	for result.Next(ctx) {

		item := reflect.New(itemTyp)
		err := result.Decode(item.Interface())
		if err != nil {
			log.Println(err)
			break
		}

		slice = reflect.Append(slice, reflect.Indirect(item))
	}
	val.Elem().Set(slice)
}

func (collection *collection) Delete2() (count int64, err error) {
	if collection.filter == nil || len(collection.filter) == 0 {
		log.Println("you can't delete all documents, it's very dangerous")
		err = errors.New("you can't delete all documents, it's very dangerous")
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.DeleteMany(ctx, collection.filter)
	if err != nil {
		log.Println(err)
		return
	}
	count = result.DeletedCount
	return
}

// 删除数据,并返回删除成功的数量
// Deprecated: 将被删除，请调用 Delete2
func (collection *collection) Delete() int64 {
	if collection.filter == nil || len(collection.filter) == 0 {
		log.Println("you can't delete all documents, it's very dangerous")
		return 0
	}
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.DeleteMany(ctx, collection.filter)
	if err != nil {
		log.Println(err)
	}
	return result.DeletedCount
}

func (collection *collection) Count2() (result int64, err error) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err = collection.Table.CountDocuments(ctx, collection.filter)
	if err != nil {
		log.Println(err)
		return
	}
	return
}

// Deprecated: 将被删除，请调用 Count2
func (collection *collection) Count() int64 {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	result, err := collection.Table.CountDocuments(ctx, collection.filter)
	if err != nil {
		log.Println(err)
		return 0
	}
	return result
}

func BeforeCreate(document interface{}) interface{} {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)

	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeCreate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]interface{}, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			sliceData[i] = BeforeCreate(val.Index(i).Interface()).(bson.M)
		}
		return sliceData

	case reflect.Struct:
		var data = make(bson.M)
		for i := 0; i < typ.NumField(); i++ {
			data[typ.Field(i).Tag.Get("bson")] = val.Field(i).Interface()
		}
		dataVal := reflect.ValueOf(data)
		if val.FieldByName("Id").Type() == reflect.TypeOf(primitive.ObjectID{}) {
			dataVal.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID()))
		}

		if val.FieldByName("Id").Interface() == "" {
			dataVal.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID().String()))
		}

		dataVal.SetMapIndex(reflect.ValueOf("created_at"), reflect.ValueOf(time.Now().Unix()))
		dataVal.SetMapIndex(reflect.ValueOf("updated_at"), reflect.ValueOf(time.Now().Unix()))
		return dataVal.Interface()

	default:
		if val.Type() == reflect.TypeOf(bson.M{}) {
			if !val.MapIndex(reflect.ValueOf("_id")).IsValid() {
				val.SetMapIndex(reflect.ValueOf("_id"), reflect.ValueOf(primitive.NewObjectID()))
			}
			val.SetMapIndex(reflect.ValueOf("created_at"), reflect.ValueOf(time.Now().Unix()))
			val.SetMapIndex(reflect.ValueOf("updated_at"), reflect.ValueOf(time.Now().Unix()))
		}
		return val.Interface()
	}
}

func BeforeUpdate(document interface{}) interface{} {
	val := reflect.ValueOf(document)
	typ := reflect.TypeOf(document)

	switch typ.Kind() {
	case reflect.Ptr:
		return BeforeUpdate(val.Elem().Interface())

	case reflect.Array, reflect.Slice:
		var sliceData = make([]interface{}, val.Len(), val.Cap())
		for i := 0; i < val.Len(); i++ {
			sliceData[i] = BeforeCreate(val.Index(i).Interface()).(bson.M)
		}
		return sliceData

	case reflect.Struct:
		var data = make(bson.M)
		for i := 0; i < typ.NumField(); i++ {
			_, ok := typ.Field(i).Tag.Lookup("over")
			if ok {
				continue
			}
			data[typ.Field(i).Tag.Get("bson")] = val.Field(i).Interface()
		}
		dataVal := reflect.ValueOf(data)
		dataVal.SetMapIndex(reflect.ValueOf("updated_at"), reflect.ValueOf(time.Now().Unix()))
		return dataVal.Interface()

	default:
		if val.Type() == reflect.TypeOf(bson.M{}) {
			val.SetMapIndex(reflect.ValueOf("updated_at"), reflect.ValueOf(time.Now().Unix()))
		}
		return val.Interface()
	}
}

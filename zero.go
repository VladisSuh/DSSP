package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type StringPoolManager struct {
	Pools map[string]string
	mu    sync.Mutex
}

var (
	instance *StringPoolManager
	once     sync.Once
)

func GetStringPoolManager() *StringPoolManager {
	once.Do(func() {
		instance = &StringPoolManager{
			Pools: make(map[string]string),
		}
	})
	return instance
}

func (sp *StringPoolManager) Get(str string) string {
	sp.mu.Lock()
	defer sp.mu.Unlock()

	if s, exists := sp.Pools[str]; exists {
		return s
	}
	sp.Pools[str] = str
	return str
}

type Tree interface {
	Insert(key string, value interface{}) error
	Get(key string) (interface{}, error)
	GetRange(minValue, maxValue string) ([]string, error)
	Update(key string, value interface{}) error
	Remove(key string) error
	SaveToFile(filename string) error
}

type TreeManager struct {
	Tree Tree
}

func NewTreeManager(treeType string) *TreeManager {
	var tree Tree
	switch treeType {
	case "avl":
		tree = NewAVLTree()
	case "redblack":
		tree = NewRedBlackTree()
	case "btree":
		tree = NewBTree()
	default:
		tree = NewMapCollection()
	}
	return &TreeManager{Tree: tree}
}

func (tc *TreeManager) Insert(key string, value interface{}) error {
	return tc.Tree.Insert(key, value)
}

func (tc *TreeManager) Get(key string) (interface{}, error) {
	return tc.Tree.Get(key)
}

func (tc *TreeManager) GetRange(minValue, maxValue string) ([]string, error) {
	return tc.Tree.GetRange(minValue, maxValue)
}

func (tc *TreeManager) Update(key string, value interface{}) error {
	return tc.Tree.Update(key, value)
}

func (tc *TreeManager) Remove(key string) error {
	return tc.Tree.Remove(key)
}

func (tc *TreeManager) SaveToFile(filename string) error {
	return tc.Tree.SaveToFile(filename)
}

type MapCollection struct {
	Data map[string]interface{}
}

func NewMapCollection() *MapCollection {
	return &MapCollection{
		Data: make(map[string]interface{}),
	}
}

func (mc *MapCollection) Insert(key string, value interface{}) error {
	sp := GetStringPoolManager()
	key = sp.Get(key)

	if _, exists := mc.Data[key]; exists {
		return errors.New("Элемент с таким ключом уже существует!")
	}
	mc.Data[key] = value
	fmt.Println("Элемент успешно добавлен с ключом", key)
	return nil
}

func (mc *MapCollection) Get(key string) (interface{}, error) {
	sp := GetStringPoolManager()
	key = sp.Get(key)

	value, exists := mc.Data[key]
	if !exists {
		return nil, errors.New("Элемент не найден!")
	}
	return value, nil
}

func (mc *MapCollection) GetRange(minValue, maxValue string) ([]string, error) {
	sp := GetStringPoolManager()
	minValue = sp.Get(minValue)
	maxValue = sp.Get(maxValue)

	var result []string
	for key := range mc.Data {
		if key >= minValue && key <= maxValue {
			result = append(result, key)
		}
	}
	return result, nil
}

func (mc *MapCollection) Update(key string, value interface{}) error {
	sp := GetStringPoolManager()
	key = sp.Get(key)

	if _, exists := mc.Data[key]; !exists {
		return errors.New("Элемент не найден!")
	}
	mc.Data[key] = value
	fmt.Println("Значение элемента с ключом", key, "успешно обновлено.")
	return nil
}

func (mc *MapCollection) Remove(key string) error {
	sp := GetStringPoolManager()
	key = sp.Get(key)

	if _, exists := mc.Data[key]; !exists {
		return errors.New("Элемент не найден!")
	}
	delete(mc.Data, key)
	return nil
}

func (mc *MapCollection) SaveToFile(filename string) error {
	data, err := json.Marshal(mc.Data)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

type PoolManager struct {
	Pools map[string]*Pool
}

func NewPoolManager() *PoolManager {
	return &PoolManager{
		Pools: make(map[string]*Pool),
	}
}

func (pm *PoolManager) ShowAll() {
	fmt.Println("Текущие пулы и схемы:")
	for poolName, pool := range pm.Pools {
		fmt.Printf("Пул: %s\n", poolName)
		for schemaName := range pool.Schemas {
			fmt.Printf("  Схема: %s\n", schemaName)
		}
	}
}

func (pm *PoolManager) AddPool(name string) {
	if _, exists := pm.Pools[name]; exists {
		fmt.Println("Пул с именем", name, "уже существует.")
	} else {
		pm.Pools[name] = NewPool()
		fmt.Println("Добавлен пул с именем", name)
	}
	pm.ShowAll()
}

func (pm *PoolManager) RemovePool(name string) {
	if pool, exists := pm.Pools[name]; exists {
		for schemaName := range pool.Schemas {
			schema := pool.Schemas[schemaName]
			for collectionName := range schema.Collections {
				schema.RemoveCollection(collectionName)
			}
			pool.RemoveSchema(schemaName)
		}
		delete(pm.Pools, name)
		fmt.Println("Пул с именем", name, "удален.")
	} else {
		fmt.Println("Пул с именем", name, "не существует.")
	}
	pm.ShowAll()
}

func (pm *PoolManager) GetPool(name string) (*Pool, error) {
	pool, ok := pm.Pools[name]
	if !ok {
		return nil, errors.New("Элемент не найден!")
	}
	return pool, nil
}

func (pm *PoolManager) GetRange(minValue, maxValue string) ([]*Pool, error) {
	var result []*Pool
	for name, pool := range pm.Pools {
		if name >= minValue && name <= maxValue {
			result = append(result, pool)
		}
	}
	return result, nil
}

func (pm *PoolManager) SaveToFile(filename string) error {
	data, err := json.Marshal(pm)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

type Pool struct {
	Schemas map[string]*Schema
}

func NewPool() *Pool {
	return &Pool{
		Schemas: make(map[string]*Schema),
	}
}

func (p *Pool) GetSchema(schemaName string) (*Schema, error) {
	schema, ok := p.Schemas[schemaName]
	if !ok {
		return nil, errors.New("Элемент не найден!")
	}
	return schema, nil
}

func (p *Pool) AddSchema(name string) {
	if _, exists := p.Schemas[name]; exists {
		fmt.Println("Схема с именем", name, "уже существует в пуле.")
	} else {
		p.Schemas[name] = NewSchema()
		fmt.Println("Схема с именем", name, "добавлена в пул.")
	}
	p.ShowSchemas()
}

func (p *Pool) RemoveSchema(name string) {
	if schema, exists := p.Schemas[name]; exists {
		for collectionName := range schema.Collections {
			schema.RemoveCollection(collectionName)
		}
		delete(p.Schemas, name)
		fmt.Println("Схема с именем", name, "удалена из пула.")
	} else {
		fmt.Println("Схема с именем", name, "не найдена в пуле.")
	}
	p.ShowSchemas()
}

func (p *Pool) ShowSchemas() {
	fmt.Println("Текущие схемы в пуле:")
	for schemaName := range p.Schemas {
		fmt.Printf("  Схема: %s\n", schemaName)
	}
}

func (p *Pool) SaveToFile(filename string) error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

type Schema struct {
	Collections map[string]TreeManager
}

func NewSchema() *Schema {
	return &Schema{
		Collections: make(map[string]TreeManager),
	}
}

func (s *Schema) GetCollection(name string) (TreeManager, error) {
	collection, ok := s.Collections[name]
	if !ok {
		return TreeManager{}, errors.New("Элемент не найден!")
	}
	return collection, nil
}

func (p *Pool) AddCollection(schemaName, collectionName string, collection TreeManager) error {
	schema, err := p.GetSchema(schemaName)
	if err != nil {
		return err
	}

	if _, exists := schema.Collections[collectionName]; exists {
		return errors.New("Коллекция с таким именем уже существует!")
	}

	schema.Collections[collectionName] = collection
	fmt.Printf("Коллекция с именем %s добавлена в схему %s в пуле\n", collectionName, schemaName)
	schema.ShowCollections()
	return nil
}

func (s *Schema) RemoveCollection(name string) {
	if _, exists := s.Collections[name]; exists {
		delete(s.Collections, name)
		fmt.Println("Коллекция с именем", name, "удалена из схемы.")
	} else {
		fmt.Println("Коллекция с именем", name, "не найдена в схеме.")
	}
	s.ShowCollections()
}

func (s *Schema) ShowCollections() {
	fmt.Println("Текущие коллекции в схеме:")
	for collectionName := range s.Collections {
		fmt.Printf("  Коллекция: %s\n", collectionName)
	}
}

func (s *Schema) SaveToFile(filename string) error {
	data, err := json.Marshal(s)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

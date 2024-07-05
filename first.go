package main

import (
	"fmt"
	"strings"
	"time"
)

func handlePoolsAndSchemas(pools *PoolManager, args []string) error {
	pools.ShowAll()
	if len(args) < 2 {
		return fmt.Errorf("недостаточно аргументов для команды %s", args[0])
	}

	switch args[0] {
	case "add-pool":
		pools.AddPool(args[1])
	case "remove-pool":
		pools.RemovePool(args[1])
	case "add-schema":
		if len(args) < 3 {
			return fmt.Errorf("недостаточно аргументов для команды add-schema")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		pool.AddSchema(args[2])
	case "remove-schema":
		if len(args) < 3 {
			return fmt.Errorf("недостаточно аргументов для команды remove-schema")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		pool.RemoveSchema(args[2])
	case "add-collection":
		if len(args) < 5 {
			return fmt.Errorf("недостаточно аргументов для команды add-collection")
		}
		collectionType := args[4]
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		treeCollection := NewTreeManager(collectionType)
		if err = pool.AddCollection(args[2], args[3], *treeCollection); err != nil {
			return err
		}
	case "remove-collection":
		if len(args) < 4 {
			return fmt.Errorf("недостаточно аргументов для команды remove-collection")
		}
		pool, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		schema, err := pool.GetSchema(args[2])
		if err != nil {
			return err
		}
		schema.RemoveCollection(args[3])
	default:
		return fmt.Errorf("неизвестная команда")
	}
	return nil
}

func runCommand(pools *PoolManager, command string, cr *ChainOfResponsibility) error {
	args := strings.Fields(command)
	if len(args) == 0 {
		return fmt.Errorf("не указана команда")
	}

	switch args[0] {
	case "add-pool", "remove-pool", "add-schema", "remove-schema", "add-collection", "remove-collection":
		return handlePoolsAndSchemas(pools, args)
	case "insert-data":
		if len(args) < 5 {
			return fmt.Errorf("недостаточно аргументов для команды insert-data")
		}
		data := TData{Key: args[3], Value: args[4]}
		insertCmd := &InsertCommand{InitialVersion: data}
		cr.AddHandler(insertCmd)
		fmt.Println("Команда вставки добавлена")
	case "update-data":
		if len(args) < 4 {
			return fmt.Errorf("недостаточно аргументов для команды update-data")
		}
		updateCmd := &UpdateCommand{UpdateExpression: args[3]}
		cr.AddHandler(updateCmd)
		fmt.Println("Команда обновления добавлена")
	case "delete-data":
		if len(args) < 3 {
			return fmt.Errorf("недостаточно аргументов для команды delete-data")
		}
		deleteCmd := &DisposeCommand{}
		cr.AddHandler(deleteCmd)
		fmt.Println("Команда удаления добавлена")
	case "get-data":
		if len(args) < 2 {
			return fmt.Errorf("недостаточно аргументов для команды get-data")
		}
		data, err := pools.GetPool(args[1])
		if err != nil {
			return err
		}
		fmt.Println("Полученные данные:", data)
	case "execute":
		var dataExists bool
		var data TData
		data.Timestamp = time.Now()
		cr.FirstHandler.Handle(&dataExists, &data, time.Now().Unix())
		handleCommand(&data)
		fmt.Println("Команды выполнены, текущее состояние:", data.Timestamp.Format("2006-01-02 15:04:05"))
	case "save-state":
		if len(args) < 2 {
			return fmt.Errorf("недостаточно аргументов для команды save-state")
		}
		err := pools.SaveToFile(args[1])
		if err != nil {
			return err
		}
		fmt.Println("Состояние системы успешно сохранено в файл:", args[1])
	case "exit":
		return nil
	default:
		return fmt.Errorf("неизвестная команда")
	}
	return nil
}

func handleCommand(data *TData) {
	data.Timestamp = time.Now()
}

var users = map[string]string{
	"admin": "password1234",
}

func authenticate(username, password string) bool {
	if pass, ok := users[username]; ok {
		return pass == password
	}
	return false
}

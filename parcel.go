package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("INSERT INTO parcel (client, status, address, created_at) VALUES (?, ?, ?, ?)", p.Client, p.Status, p.Address, p.CreatedAt)
	// верните идентификатор последней добавленной записи
	if err != nil {
		return 0, fmt.Errorf("не удалось создать запись в базе данных: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("не удалось найти идентификатор последней добавленной записи: %w", err)
	}
	return int(id), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	row, err := s.db.Query("SELECT * FROM parcel WHERE number = ?", number)
	// заполните объект Parcel данными из таблицы
	if err != nil {
		return Parcel{}, fmt.Errorf("не удалось получить запись по заданному number: %w", err)
	}
	p := Parcel{}
	error := row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if error != nil {
		return p, fmt.Errorf("не удалось получить данные из строки: %w", err)
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	rows, err := s.db.Query("SELECT * FROM parcel WHERE client = ?", client)
	// заполните срез Parcel данными из таблицы

	if err != nil {
		return nil, fmt.Errorf("не удалось получить записи по заданному client: %w", err)
	}
	var res []Parcel
	for rows.Next() {
		p := Parcel{}
		err := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("не записать данные из строки: %w", err)
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("UPDATE parcel SET status = ? WHERE number = ?", status, number)
	if err != nil {
		return fmt.Errorf("не удалось обновить данные: %w", err)
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	_, err := s.db.Exec("UPDATE parcel SET address = ? WHERE number = ?", address, number)
	if err != nil {
		return fmt.Errorf("не удалось обновить данные: %w", err)
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	_, err := s.db.Exec("DELETE FROM parcel WHERE number = ? AND status = ?", number, ParcelStatusRegistered)

	if err != nil {
		return fmt.Errorf("не удалось удалить строку: %w", err)
	}
	return nil
}

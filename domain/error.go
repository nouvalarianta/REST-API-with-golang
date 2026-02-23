package domain

import "errors"

var BookNotFound = errors.New("data buku tidak di temukan")
var JournalNotFound = errors.New("journal buku tidak ditemukan")

var JournalAlreadyCompleted = errors.New("buku sudah di kembalikan")

var CustomerNotFound = errors.New("data customer tidak di temukan")
var InvalidID = errors.New("id tidak valid")

var BookAlreadyBorrowed = errors.New("buku sudah di pinjam")
var BookStockNotFound = errors.New("stok buku tidak di temukan")
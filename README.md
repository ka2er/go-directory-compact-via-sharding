# go-directory-compact-via-sharding
Simple utility to shard a directory

## Purpose

I need a tool to prepare music USB stick for my Idrive enabled car as scrolling is slow... and wanted to take a look at golang.
So I wrote this piece of code .

## Install

Just type the following
```
go install
```

## Usage

You just need to pass from and destination directory.
 - Destination is created if needed
 - Max parameter is optionnal

```
Usage of go-directory-compact-via-sharding:
  -from="": Directory to process - mandatory
  -max=10: Max number of TOP directory - optionnal
  -to="": Destination directory - mandatory
```

## Exemple

Before 
```
├── antipapalist
│   └── pentadecahydrated
├── ignobly
│   └── purplely
├── induction
│   └── tetralemma
├── jemadar
│   └── panharmonicon
├── minatorially
│   └── embryotrophy
├── noseless
│   └── stereotelemeter
├── platyhelminthes
│   └── imbased
├── puerile
│   └── adital
└── semicontinuum
    └── frontiersman

```
After
```
├── A => I
│   ├── antipapalist
│   ├── ignobly
│   └── induction
├── J => N
│   ├── jemadar
│   ├── minatorially
│   └── noseless
└── P => S
    ├── platyhelminthes
    ├── puerile
    └── semicontinuum
```    

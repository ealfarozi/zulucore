# zulucore
## github information
Flow pada github: https://guides.github.com/introduction/flow/

## Cara pull ke local
Buka CMD/Terminal, CD ke folder project yg masih kosong
`git clone https://github.com/ealfarozi/zulucore.git`

## cara commit
commit adalah flag files yang sudah dirubah dan di upload masih di local (laptop)
- `git init` untuk menandakan bahwa folder tersebut adalah folder yang akan di push ke github
- `git add .` "." bisa diganti dengan nama file
- `git commit -m "Message"` -m adalah message atau catatan untuk commit
- `git remote add origin https://github.com/ealfarozi/zulucore.git` untuk flag tujuan, disini harus masukkan username dan password github
- `git push -u origin development` development adalah branch selama masa development
- jika ingin push ke master/production maka:
`git push -u origin master`

PS:
- jika ingin pindah-pindah branch untuk di push bisa menggunakan perintah `git branch` untuk lihat list branch yang sedang aktif, yang sedang aktif mudah terlihat dari tanda asterisk/bintang.
- jika ingin dirubah branch nya bisa menggunakan perintah `git checkout master` untuk pindah ke branch master atau `git checkout development` untuk pindah ke branch development. setelah itu baru jalan perintah `git push -u origin <branch_name>`

## cara pull
cara pull cukup mudah, hanya `git pull`

## how to start
- `go mod init` --> init go mod file
- `go mod tidy` --> listing all of the dependencies
- create your own `.env` file (set your JWT_KEY and DB_URL)
- DB_URL format `user:pwd@(url:port)/dbname`, example: `root:password@(localhost:3306)/zulu`
- remove go.mod and go.sum files
- `go run main.go` after this command runs, go will downloading all of the packages

# porm
gin-pro 的 orm

## 提供一个更加简洁的orm，用于快速开发小型项目

## 快速入门

`go get github.com/gin-pro/porm `

```
func TestGet(t *testing.T) {
	db, err := NewEngine("mysql", "xxx")
	if err != nil {
		fmt.Println(fmt.Sprintf("open err : %v", err))
		return
	}
	type User struct {
		ID   int    `porm:"id"`
		Name string `porm:"name"`
		Sex  string `porm:"-"`
	}
	user := *&User{}
	_, err = db.Table("user").Where("id = ?", 2).Get(&user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Get err : %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("User : %v", user))
}

func TestFind(t *testing.T) {
	db, err := NewEngine("mysql", "xx")
	if err != nil {
		fmt.Println(fmt.Sprintf("open err : %v", err))
		return
	}
	type User struct {
		ID   int    `porm:"id"`
		Name string `porm:"name"`
		Sex  string `porm:"-"`
	}
	user := []User{}
	_, err = db.Table("user").Find(&user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Get err : %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("User : %v", user))
}
  
  
func TestInsert(t *testing.T) {
	db, err := NewEngine("mysql", "xxx")
	if err != nil {
		fmt.Println(fmt.Sprintf("open err : %v", err))
		return
	}
	type User struct {
		ID   int    `porm:"id"`
		Name string `porm:"name"`
		Sex  string `porm:"-"`
	}
	user := &User{0, "ttt", "123"}
	_, err = db.NewSession().Table("user").Insert(user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Get err : %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("User : %v", user))
}

func TestUpdate(t *testing.T) {
	db, err := NewEngine("mysql", "xxx")
	if err != nil {
		fmt.Println(fmt.Sprintf("open err : %v", err))
		return
	}
	type User struct {
		Name string `porm:"name"`
		Sex  string `porm:"-"`
	}
	user := &User{"ttt1", "123"}
	_, err = db.NewSession().Table("user").Where("id = ?", 2).Update(user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Get err : %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("User : %v", user))
}


func TestDelete(t *testing.T) {
	db, err := NewEngine("mysql", "xxx")
	if err != nil {
		fmt.Println(fmt.Sprintf("open err : %v", err))
		return
	}
	type User struct {
		Name string `porm:"name"`
		Sex  string `porm:"-"`
	}
	user := &User{"ttt1", "123"}
	_, err = db.NewSession().Table("user").Where("id = ?", 2).Delete(user)
	if err != nil {
		fmt.Println(fmt.Sprintf("Get err : %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("User : %v", user))
}
```
## 开发日志

日期 ：2022-01-09

内容：

```
测试版本1.0

- 提供增删查改四个方法

- 提供打开连接
```

input := "$5\r\nAhmed\r\n"
reader := bufio.NewReader(strings.NewReader(input))

b, _ := reader.ReadBytes()
if b != '$'{
	fmt.Println("Invalid type, expecting bulk strings only")
    os.Exit(1)
}

size 
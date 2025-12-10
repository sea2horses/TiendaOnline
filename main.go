package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"tienda-online/internal"
	"tienda-online/internal/db"
	"tienda-online/models"

	_ "github.com/microsoft/go-mssqldb"
)

const (
	defaultServer = "localhost"
	defaultDB     = "Tienda"
)

const (
	colorReset   = "\033[0m"
	colorCyan    = "\033[36m"
	colorGreen   = "\033[32m"
	colorMagenta = "\033[35m"
	colorYellow  = "\033[33m"
	colorRed     = "\033[31m"
)

var scanner = bufio.NewScanner(os.Stdin)

func main() {
	printBanner()

	conn, err := initDatabaseFlow(defaultServer)
	internal.Check(err, "No se pudo inicializar la base de datos")
	defer conn.Close()
	db.SetDatabase(conn)

	mainMenu()
	fmt.Println("Hasta luego")
}

func printBanner() {
	fmt.Println(string(colorCyan) + "╔══════════════════════════════════════════════╗")
	fmt.Println("║            TIENDA ONLINE CONSOLE             ║")
	fmt.Println("╠══════════════════════════════════════════════╣")
	fmt.Println("║   CRUD interactivo para todos los modelos    ║")
	fmt.Println("╚══════════════════════════════════════════════╝" + string(colorReset))
}

func initDatabaseFlow(server string) (*sql.DB, error) {
	fmt.Printf("%sIntentando conectar a %s/%s...%s\n", colorYellow, server, defaultDB, colorReset)
	conn, err := db.AttemptConnection(server, defaultDB)
	if err == nil {
		fmt.Printf("%sConectado a %s/%s%s\n", colorGreen, server, defaultDB, colorReset)
		return conn, nil
	}

	fmt.Printf("%sNo se pudo conectar a %s/%s: %v%s\n", colorRed, server, defaultDB, err, colorReset)
	fmt.Printf("%sCrearemos la base y tablas con queries/init.sql%s\n", colorYellow, colorReset)

	masterConn, err := db.AttemptConnection(server, "master")
	if err != nil {
		return nil, fmt.Errorf("no se pudo conectar a master: %w", err)
	}
	db.SetDatabase(masterConn)
	fmt.Println("Conectado a master.")
	_, err = db.ExecFromFile(context.Background(), "init.sql")
	masterConn.Close()
	if err != nil {
		return nil, fmt.Errorf("fallo ejecutando init.sql: %w", err)
	}

	conn, err = db.AttemptConnection(server, defaultDB)
	if err != nil {
		return nil, fmt.Errorf("init.sql corrio pero no se pudo conectar a %s: %w", defaultDB, err)
	}
	fmt.Printf("%sBase de datos %s creada y conectada%s\n", colorGreen, defaultDB, colorReset)

	if confirm("¿Deseas cargar datos de prueba (queries/init_data.sql)? (s/N): ") {
		db.SetDatabase(conn)
		if _, err := db.ExecFromFile(context.Background(), "init_data.sql"); err != nil {
			return nil, fmt.Errorf("fallo insertando datos de prueba: %w", err)
		}
		fmt.Printf("%sDatos de prueba insertados.%s\n", colorGreen, colorReset)
	}

	return conn, nil
}

func mainMenu() {
	for {
		fmt.Println()
		fmt.Println(colorMagenta + "=== MENU PRINCIPAL ===" + colorReset)
		fmt.Println("[1] Clientes")
		fmt.Println("[2] Categorias")
		fmt.Println("[3] Productos")
		fmt.Println("[4] SKUs")
		fmt.Println("[5] Carritos")
		fmt.Println("[6] Detalles de carrito")
		fmt.Println("[7] Reseñas")
		fmt.Println("[8] Direcciones")
		fmt.Println("[9] Pedidos")
		fmt.Println("[10] Devoluciones")
		fmt.Println("[I] Re-ejecutar init.sql")
		fmt.Println("[D] Insertar datos de prueba (init_data.sql)")
		fmt.Println("[Q] Salir")

		choice := readLine("Elige una opcion: ")
		switch strings.ToLower(choice) {
		case "1":
			menuClientes()
		case "2":
			menuCategorias()
		case "3":
			menuProductos()
		case "4":
			menuSKUs()
		case "5":
			menuCarritos()
		case "6":
			menuDetalles()
		case "7":
			menuResenas()
		case "8":
			menuDirecciones()
		case "9":
			menuPedidos()
		case "10":
			menuDevoluciones()
		case "i":
			runInit()
		case "d":
			runInitData()
		case "q":
			return
		default:
			fmt.Println(colorRed + "Opcion no valida" + colorReset)
		}
	}
}

func runInit() {
	db.SetDatabase(db.CurrentDatabase)
	if _, err := db.ExecFromFile(context.Background(), "init.sql"); err != nil {
		fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
	} else {
		fmt.Printf("%sSchema reinicializado.%s\n", colorGreen, colorReset)
	}
}

func runInitData() {
	db.SetDatabase(db.CurrentDatabase)
	if _, err := db.ExecFromFile(context.Background(), "init_data.sql"); err != nil {
		fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
	} else {
		fmt.Printf("%sDatos de prueba insertados.%s\n", colorGreen, colorReset)
	}
}

// ===== Menus por entidad =====

func menuClientes() {
	m := models.NewClienteManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- Clientes --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Cambiar contraseña")
		fmt.Println("[6] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			name := readLine("Nombre: ")
			phone := readLine("Telefono (8 digitos): ")
			email := readLine("Correo (opcional): ")
			pass := readLine("Contraseña: ")
			handleErr(m.Create(context.Background(), name, phone, email, pass))
		case "4":
			id := readInt("ID: ")
			name := readLine("Nombre: ")
			phone := readLine("Telefono: ")
			email := readLine("Correo (opcional): ")
			handleErr(m.Update(context.Background(), id, name, phone, email))
		case "5":
			id := readInt("ID: ")
			pass := readLine("Nueva contraseña: ")
			handleErr(m.UpdatePassword(context.Background(), id, pass))
		case "6":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

func menuCategorias() {
	m := models.NewCategoriaManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- Categorias --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			name := readLine("Nombre: ")
			handleErr(m.Create(context.Background(), name))
		case "4":
			id := readInt("ID: ")
			name := readLine("Nombre: ")
			handleErr(m.Update(context.Background(), id, name))
		case "5":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

func menuProductos() {
	m := models.NewProductoManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- Productos --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			desc := readLine("Descripcion: ")
			catID := readOptionalInt("ID Categoria (0 para NULL): ")
			handleErr(m.Create(context.Background(), desc, catID))
		case "4":
			id := readInt("ID: ")
			desc := readLine("Descripcion: ")
			catID := readOptionalInt("ID Categoria (0 para NULL): ")
			handleErr(m.Update(context.Background(), id, desc, catID))
		case "5":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

func menuSKUs() {
	m := models.NewSKUManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- SKUs --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			pid := readInt("ID Producto: ")
			price := readFloat("Precio: ")
			stock := readInt("Stock: ")
			handleErr(m.Create(context.Background(), pid, price, stock))
		case "4":
			id := readInt("ID: ")
			pid := readInt("ID Producto: ")
			price := readFloat("Precio: ")
			stock := readInt("Stock: ")
			handleErr(m.Update(context.Background(), id, pid, price, stock))
		case "5":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

func menuCarritos() {
	m := models.NewCarritoManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- Carritos --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			uid := readInt("ID Usuario: ")
			handleErr(m.Create(context.Background(), uid))
		case "4":
			id := readInt("ID: ")
			uid := readInt("ID Usuario: ")
			handleErr(m.Update(context.Background(), id, uid))
		case "5":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

func menuDetalles() {
	m := models.NewCarritoDetalleManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- Detalles de Carrito --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			cid := readInt("ID Carrito: ")
			sid := readInt("ID SKU: ")
			qty := readInt("Cantidad: ")
			handleErr(m.Create(context.Background(), cid, sid, qty))
		case "4":
			id := readInt("ID: ")
			cid := readInt("ID Carrito: ")
			sid := readInt("ID SKU: ")
			qty := readInt("Cantidad: ")
			handleErr(m.Update(context.Background(), id, cid, sid, qty))
		case "5":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

func menuResenas() {
	m := models.NewResenaManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- Reseñas --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			uid := readInt("ID Usuario: ")
			pid := readInt("ID Producto: ")
			r := readInt("Puntuacion (1-5): ")
			comment := readLine("Comentario (opcional): ")
			handleErr(m.Create(context.Background(), uid, pid, r, comment))
		case "4":
			id := readInt("ID: ")
			uid := readInt("ID Usuario: ")
			pid := readInt("ID Producto: ")
			r := readInt("Puntuacion (1-5): ")
			comment := readLine("Comentario (opcional): ")
			handleErr(m.Update(context.Background(), id, uid, pid, r, comment))
		case "5":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

func menuDirecciones() {
	m := models.NewDireccionManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- Direcciones --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			uid := readInt("ID Usuario: ")
			tipo := readLine("Tipo (Envio/Facturacion): ")
			detalle := readLine("Detalle (opcional): ")
			handleErr(m.Create(context.Background(), uid, tipo, detalle))
		case "4":
			id := readInt("ID: ")
			uid := readInt("ID Usuario: ")
			tipo := readLine("Tipo (Envio/Facturacion): ")
			detalle := readLine("Detalle (opcional): ")
			handleErr(m.Update(context.Background(), id, uid, tipo, detalle))
		case "5":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

func menuPedidos() {
	m := models.NewPedidoManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- Pedidos --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			uid := readInt("ID Usuario: ")
			del := confirm("¿Entregado? (s/N): ")
			handleErr(m.Create(context.Background(), uid, del))
		case "4":
			id := readInt("ID: ")
			uid := readInt("ID Usuario: ")
			del := confirm("¿Entregado? (s/N): ")
			handleErr(m.Update(context.Background(), id, uid, del))
		case "5":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

func menuDevoluciones() {
	m := models.NewDevolucionManager(db.CurrentDatabase)
	for {
		fmt.Println(colorCyan + "\n-- Devoluciones --" + colorReset)
		fmt.Println("[1] Listar")
		fmt.Println("[2] Ver por ID")
		fmt.Println("[3] Crear")
		fmt.Println("[4] Actualizar")
		fmt.Println("[5] Eliminar")
		fmt.Println("[B] Volver")
		c := readLine("Opcion: ")
		switch strings.ToLower(c) {
		case "1":
			items, err := m.List(context.Background())
			if handleErr(err) {
				break
			}
			internal.ListItems(items)
		case "2":
			id := readInt("ID: ")
			item, err := m.Get(context.Background(), id)
			if handleErr(err) {
				break
			}
			fmt.Printf("%+v\n", item)
		case "3":
			oid := readInt("ID Pedido: ")
			fecha := readDate("Fecha (YYYY-MM-DD): ")
			estado := readLine("Estado (opcional): ")
			desc := readLine("Descripcion (opcional): ")
			res := readLine("Resolucion (opcional): ")
			handleErr(m.Create(context.Background(), oid, fecha, estado, desc, res))
		case "4":
			id := readInt("ID: ")
			oid := readInt("ID Pedido: ")
			fecha := readDate("Fecha (YYYY-MM-DD): ")
			estado := readLine("Estado (opcional): ")
			desc := readLine("Descripcion (opcional): ")
			res := readLine("Resolucion (opcional): ")
			handleErr(m.Update(context.Background(), id, oid, fecha, estado, desc, res))
		case "5":
			id := readInt("ID: ")
			if confirm("¿Seguro? (s/N): ") {
				handleErr(m.Delete(context.Background(), id))
			}
		case "b":
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

// ===== Helpers de entrada =====

func readLine(prompt string) string {
	fmt.Print(prompt)
	for {
		if scanner.Scan() {
			return strings.TrimSpace(scanner.Text())
		}
		// Si fallo (por buffer u otro motivo), re-creamos el scanner y reintentamos.
		scanner = bufio.NewScanner(os.Stdin)
		fmt.Print(prompt)
	}
}

func readInt(prompt string) int {
	for {
		val := readLine(prompt)
		i, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Ingresa un numero valido")
			continue
		}
		return i
	}
}

func readOptionalInt(prompt string) int {
	for {
		val := readLine(prompt)
		if strings.TrimSpace(val) == "" {
			return 0
		}
		i, err := strconv.Atoi(val)
		if err != nil {
			fmt.Println("Ingresa un numero valido")
			continue
		}
		return i
	}
}

func readFloat(prompt string) float64 {
	for {
		val := readLine(prompt)
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			fmt.Println("Ingresa un numero valido")
			continue
		}
		return f
	}
}

func readDate(prompt string) time.Time {
	for {
		val := readLine(prompt)
		t, err := time.Parse("2006-01-02", val)
		if err != nil {
			fmt.Println("Formato invalido, usa YYYY-MM-DD")
			continue
		}
		return t
	}
}

func confirm(prompt string) bool {
	val := strings.ToLower(readLine(prompt))
	return val == "s" || val == "si" || val == "sí"
}

func handleErr(err error) bool {
	if err != nil {
		fmt.Printf("%sError: %v%s\n", colorRed, err, colorReset)
		return true
	}
	fmt.Printf("%sOK%s\n", colorGreen, colorReset)
	return false
}

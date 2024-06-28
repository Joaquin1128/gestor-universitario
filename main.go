package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	bolt "go.etcd.io/bbolt"
	"github.com/google/uuid"

)

type Alumno struct {
	Id_alumno int
	Nombre string
	Apellido string
	Dni int
	Fecha_nacimiento string
	Telefono string
	Email string
}

type Materia struct {
	Id_materia int
	Nombre string
}

type Comision struct {
	Id_materia int
	Id_comision int
	Cupo int
}

type Correlatividad struct {
	Id_materia int
	Id_mat_correlativa int
}

type Cursada struct {
	Id_materia int
	Id_alumno int
	Id_comision int
	F_inscripcion string
	Nota int
	Estado string
}

type Periodo struct {
	Semestre string
	Estado string
}

type Historia_academica struct {
	Id_alumno int
	Semestre string
	Id_materia int
	Id_comision int
	Estado string
	Nota_regular int
	Nota_final int
}

type Error struct {
	Id_error int
	Operacion string
	Semestre string
	Id_alumno int
	Id_materia int
	Id_comision int
	F_error string
	Motivo string
}

type Envio_email struct {
	Id_email int
	F_generacion string
	Email_alumno string
	Asunto string
	Cuerpo string
	F_envio string
	Estado string
}

type Entrada struct {
	Id_orden int
	Operacion string
	Año int
	Nro_semestre int
	Id_alumno int
	Id_materia int
	Id_comision int
	Nota int
}

func main() {
    ejecutarPrograma()
}

func mostrarOpciones() int {
	fmt.Printf ("Elige una opcion para continuar:\n")
	fmt.Printf ("Para crear la DB, escriba el nùmero 1\n")
	fmt.Printf ("Para crear las tablas de la DB, escriba el nùmero 2\n")
	fmt.Printf ("Para cargar los datos de los archivos JSON, escriba el nùmero 3\n")
	fmt.Printf ("Para agregar las primary keys, escriba el nùmero 4\n")
	fmt.Printf ("Para agregar las foreign keys, escriba el nùmero 5\n")
	fmt.Printf ("Para borrar las Primary Keys y las Foreign Keys, escriba el nùmero 6\n")
	fmt.Printf ("Para cargar todos los Stored Procedures y los Triggers, escriba el nùmero 7\n")
	fmt.Printf ("Para crear la DB y cargar todo, presione 8\n")
	fmt.Printf ("Para ejecutar las entradas, escriba el número 9\n")
	fmt.Printf ("Para crear la base de datos BoltDB y cargar sus datos, escriba el nùmero 10\n")
	fmt.Printf ("Para leer y mostrar los datos guardados en la base de datos BoltDB, escriba el nùmero 11\n")
	fmt.Printf ("Para salir, escriba el nùmero 12\n")

	var opcion int
	fmt.Scanf("%d",&opcion)
	return opcion
}

func ejecutarPrograma() {
	fmt.Printf ("¡Bienvenido!\n")
	
	connStr := "user=postgres host=localhost dbname=gestor_universitario_db1 sslmode=disable"
	
	for {
		clear()

		opcion := mostrarOpciones()

		clear()

		switch opcion {
		case 1:
		createDatabase()

		case 2:
		createDbTables()

		case 3:
		levantarJSONsSQL()

		case 4:
		agregarPrimaryKey()

		case 5:
		agregarForeignKey()

		case 6:
		borrarKeys()

		case 7:
		cargarSpTriggers(connStr)
		
		case 8:
		createDatabase()
		createDbTables()
		levantarJSONsSQL()
		agregarPrimaryKey()
		agregarForeignKey()
		cargarSpTriggers(connStr)
		
		case 9:
		ejecutarEntradas(connStr)
		
		case 10:
		levantarJSONsBoltDB()
		
		case 11:
		leerBoltDB()
		
		case 12:
		fmt.Printf("¡Hasta la proxima!\n")
		os.Exit(0)
		default:
		fmt.Printf("Opcion no valida\n")
		}

		if !preguntarContinuar() {
		break
		}
	}
}

func preguntarContinuar() bool {
	var continuar string
	fmt.Printf("Si desea realizar otra accion presione enter\n")
	fmt.Scanf("%s",&continuar)
	return continuar == ""
}

func clear() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func createDatabase() {
	db,err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err!= nil {
	log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`drop database if exists gestor_universitario_db1;`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`create database gestor_universitario_db1;`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Base de datos creada.\n")
}

func createDbTables() {
	db,err := sql.Open("postgres", "user=postgres host=localhost dbname=gestor_universitario_db1 sslmode=disable")
	if err!= nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`create table alumno(id_alumno int, nombre char(64), apellido char(64), dni int, fecha_nacimiento date, telefono char(64), email char(64));
					create table materia(id_materia int, nombre char(64));
					create table correlatividad(id_materia int, id_materia_correlativa int);
					create table comision(id_materia int, id_comision int, cupo int);
					create table cursada(id_materia int, id_alumno int, id_comision int, f_inscripcion timestamp, nota int, estado char(12));
					create table periodo(semestre char(6), estado char(15));
					create table historia_academica(id_alumno int, semestre text, id_materia int, id_comision int, estado char(15), nota_regular int, nota_final int);
					create table error(id_error int, operacion char(15), semestre text, id_alumno int, id_materia int, id_comision int, f_error timestamp, motivo char(64));
					create table envio_email(id_email int, f_generacion timestamp, email_alumno text, asunto text, cuerpo text, f_envio timestamp, estado char(10));
					create table entrada_trx(id_orden int, operacion char(15), año int, nro_semestre int, id_alumno int, id_materia int, id_comision int, nota int);`)
	if err!= nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Tablas cargadas.\n")
}

func agregarPrimaryKey (){
	db,err := sql.Open("postgres", "user=postgres host=localhost dbname=gestor_universitario_db1 sslmode=disable")
	if err!= nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`alter table alumno add constraint pk_alumno primary key (id_alumno);
					alter table materia add constraint pk_materia primary key (id_materia);
					alter table correlatividad add constraint pk_correlatividad primary key (id_materia, id_materia_correlativa);
					alter table comision add constraint pk_comision primary key (id_materia, id_comision);
					alter table cursada add constraint pk_cursada primary key (id_materia, id_alumno);
					alter table periodo add constraint pk_periodo primary key (semestre);
					alter table historia_academica add constraint pk_academica primary key (id_alumno, semestre, id_materia);
					alter table error add constraint pk_error primary key (id_error);
					alter table envio_email add constraint pk_envio_mail primary key (id_email);`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Primary Keys cargadas.\n")
}

func agregarForeignKey (){
	db,err := sql.Open("postgres", "user=postgres host=localhost dbname=gestor_universitario_db1 sslmode=disable")
	if err!= nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`alter table correlatividad add constraint fk_materia foreign key (id_materia) references materia (id_materia);
					alter table correlatividad add constraint fk_correlativa foreign key (id_materia_correlativa) references materia (id_materia);
					alter table comision add constraint fk_materia foreign key (id_materia) references materia (id_materia);
					alter table cursada add constraint fk_materia foreign key (id_materia) references materia (id_materia);
					alter table cursada add constraint fk_alumno foreign key (id_alumno) references alumno (id_alumno);
					alter table historia_academica add constraint fk_alumno foreign key (id_alumno) references alumno (id_alumno);
					alter table historia_academica add constraint fk_periodo foreign key (semestre) references periodo (semestre);
					alter table historia_academica add constraint fk_materia foreign key (id_materia) references materia (id_materia);`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Foreign Keys cargadas.\n")
}

func borrarKeys (){
	db,err := sql.Open("postgres", "user=postgres host=localhost dbname=gestor_universitario_db1 sslmode=disable")
	if err!= nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(`alter table cursada drop constraint fk_materia;
					alter table cursada drop constraint fk_alumno;
					alter table historia_academica drop constraint fk_materia;
					alter table historia_academica drop constraint fk_alumno;
					alter table historia_academica drop constraint fk_periodo;
					alter table error drop constraint fk_materia;
					alter table error drop constraint fk_alumno;
					alter table error drop constraint fk_periodo;
					alter table correlatividad drop constraint fk_materia;
					alter table correlatividad drop constraint fk_correlativa;
					alter table comision drop constraint fk_materia;
					alter table alumno drop constraint pk_alumno;
					alter table materia drop constraint pk_materia;
					alter table correlatividad drop constraint pk_correlatividad;
					alter table comision drop constraint pk_comision;
					alter table cursada drop constraint pk_cursada;
					alter table periodo drop constraint pk_periodo;
					alter table historia_academica drop constraint pk_academica;
					alter table error drop constraint pk_error;
					alter table envio_mail drop constraint pk_envio_mail;`)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Primary Keys y Foreign Keys borradas.\n")
}

func levantarJSONsSQL() {
	db,err := sql.Open("postgres", "user=postgres host=localhost dbname=gestor_universitario_db1 sslmode=disable")
	if err!= nil{
		log.Fatal(err)
	}
	defer db.Close()
	
	dataAlumnos, err := ioutil.ReadFile("data/alumnos.json")
	if err != nil{
		log.Fatal(err)
	}
	
	var alumnos []Alumno
	err = json.Unmarshal(dataAlumnos, &alumnos)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, alumno := range alumnos {
		_, err := db.Exec("insert into alumno values ($1, $2, $3, $4, $5, $6, $7)", alumno.Id_alumno, alumno.Nombre, alumno.Apellido, alumno.Dni, alumno.Fecha_nacimiento, alumno.Telefono, alumno.Email)
		if err != nil{
			log.Fatal(err)
		}
	}
	
	fmt.Printf("Tabla de alumnos cargada.\n")
	
	dataMaterias, err := ioutil.ReadFile("data/materias.json")
	if err != nil{
		log.Fatal(err)
	}
	
	var materias []Materia
	err = json.Unmarshal(dataMaterias, &materias)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, materia := range materias {
		_, err := db.Exec("insert into materia values ($1, $2)", materia.Id_materia, materia.Nombre)
		if err != nil{
			log.Fatal(err)
		}
	}
	
	fmt.Printf("Tabla de materias cargada.\n")
	
	dataComisiones, err := ioutil.ReadFile("data/comisiones.json")
	if err != nil{
		log.Fatal(err)
	}
	
	var comisiones []Comision
	err = json.Unmarshal(dataComisiones, &comisiones)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, comision := range comisiones {
		_, err := db.Exec("insert into comision values ($1, $2, $3)", comision.Id_materia, comision.Id_comision, comision.Cupo)
		if err != nil{
			log.Fatal(err)
		}
	}
	
	fmt.Printf("Tabla de comisiones cargada.\n")
	
	dataPeriodos, err := ioutil.ReadFile("data/periodos.json")
	if err != nil{
		log.Fatal(err)
	}
	
	var periodos []Periodo
	err = json.Unmarshal(dataPeriodos, &periodos)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, periodo := range periodos {
		_, err := db.Exec("insert into periodo values ($1, $2)", periodo.Semestre, periodo.Estado)
		if err != nil{
			log.Fatal(err)
		}
	}
	
	fmt.Printf("Tabla de periodos cargada.\n")
	
	dataCorrelativas, err := ioutil.ReadFile("data/correlatividades.json")
	if err != nil{
		log.Fatal(err)
	}
	
	var correlatividades []Correlatividad
	err = json.Unmarshal(dataCorrelativas, &correlatividades)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, correlativa := range correlatividades {
		_, err := db.Exec("insert into correlatividad values ($1, $2)", correlativa.Id_materia, correlativa.Id_mat_correlativa)
		if err != nil{
			log.Fatal(err)
		}
	}
	
	fmt.Printf("Tabla de correlatividades cargada.\n")
	
	dataHistorias, err := ioutil.ReadFile("data/historia_academica.json")
	if err != nil{
		log.Fatal(err)
	}
	
	var historias_academicas []Historia_academica
	err = json.Unmarshal(dataHistorias, &historias_academicas)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, historia_academica := range historias_academicas {
		_, err := db.Exec("insert into historia_academica values ($1, $2, $3, $4, $5, $6, $7)", historia_academica.Id_alumno, historia_academica.Semestre, historia_academica.Id_materia,
							historia_academica.Id_comision, historia_academica.Estado, historia_academica.Nota_regular, historia_academica.Nota_final)
		if err != nil{
			log.Fatal(err)
		}
	}
	
	fmt.Printf("Tabla de historias académicas cargada.\n")
	
	dataEntradas, err := ioutil.ReadFile("data/entradas_trx.json")
	if err != nil{
		log.Fatal(err)
	}
	
	var entradas_trx []Entrada
	err = json.Unmarshal(dataEntradas, &entradas_trx)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, entrada_trx := range entradas_trx {
		_, err := db.Exec("insert into entrada_trx values ($1, $2, $3, $4, $5, $6, $7, $8)", entrada_trx.Id_orden, entrada_trx.Operacion, entrada_trx.Año,
							entrada_trx.Nro_semestre, entrada_trx.Id_alumno, entrada_trx.Id_materia, entrada_trx.Id_comision, entrada_trx.Nota)
		if err != nil{
			log.Fatal(err)
		}
	}
	
	fmt.Printf("Tabla de entradas cargada.\n")
}

func loadSQLFilesFromFolder(connStr string, folderPath string) error {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return fmt.Errorf("Error al conectar a la base de datos: %w", err)
    }
    defer db.Close()

    files, err := ioutil.ReadDir(folderPath)
    if err != nil {
        return fmt.Errorf("Error al leer el directorio %s: %w", folderPath, err)
    }

    for _, file := range files {
        if !file.IsDir() {
            sqlFilePath := filepath.Join(folderPath, file.Name())
            err = loadSQLFile(db, sqlFilePath)
            if err != nil {
                return fmt.Errorf("Error al cargar el archivo SQL (%s): %w", sqlFilePath, err)
            }
        }
    }

    return nil
}

func crearSecuenciaEmailId() error {
	db,err := sql.Open("postgres", "user=postgres host=localhost dbname=gestor_universitario_db1 sslmode=disable")
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("create sequence if not exists envio_email_id_seq")
	return err
}

func crearSecuenciaErrorId() error {
	db,err := sql.Open("postgres", "user=postgres host=localhost dbname=gestor_universitario_db1 sslmode=disable")
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec("create sequence if not exists error_id_seq")
	return err
}

func cargarSpTriggers(connStr string){
	err := loadSQLFilesFromFolder(connStr, "stored_procedures")
	if err != nil {
		log.Fatalf("Error al cargar los Stored Procedures: %v\n", err)
	}
	fmt.Printf("Stored Procedures cargados exitosamente.\n")
	
	err = crearSecuenciaEmailId()
	if err != nil{
		log.Fatal(err)
	}
	
	err = crearSecuenciaErrorId()
	if err != nil{
		log.Fatal(err)
	}
	
	err = loadSQLFilesFromFolder(connStr, "triggers")
	if err != nil {
		log.Fatalf("Error al cargar los Triggers: %v\n", err)
	}
	fmt.Printf("Triggers cargados exitosamente.\n")
}

func ejecutarEntradas(connStr string){
	err := loadSQLFilesFromFolder(connStr, "entradas")
	if err != nil {
		log.Fatalf("Error al cargar las entradas: %v\n", err)
	}
	fmt.Printf("Entradas ejecutadas exitosamente.\n")
}

func loadSQLFile(db *sql.DB, filepath string) error {
    sql, err := ioutil.ReadFile(filepath)
    if err != nil {
        return fmt.Errorf("Error al leer el archivo SQL: %w", err)
    }

    _, err = db.Exec(string(sql))
    if err != nil {
        return fmt.Errorf("Error al ejecutar el archivo SQL: %w", err)
    }

    return nil
}

func CreateUpdate(db *bolt.DB, bucketName string, key []byte, val []byte) error {
    tx, err := db.Begin(true)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

    err = b.Put(key, val)
    if err != nil {
        return err
    }

    if err := tx.Commit(); err != nil {
        return err
    }

    return nil
}

func dropBuckets(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		bucketNames := []string{"alumno", "materia", "comision", "cursada"}

		for _, name := range bucketNames {
			if err := tx.DeleteBucket([]byte(name)); err != bolt.ErrBucketNotFound {
				if err != nil {
					return err
				}
			}
		}

		return nil
	})
}

func levantarJSONsBoltDB() {
	boltDB, err := bolt.Open("nosql.db", 0600, nil) // LA CREA SI NO EXISTE
	if err != nil {
		log.Fatal(err)
	}
	defer boltDB.Close()
	
	err = dropBuckets(boltDB)
	if err != nil {
		log.Fatal(err)
	}

	dataAlumnos, err := ioutil.ReadFile("data/alumnos.json")
	if err != nil {
		log.Fatal(err)
	}

	var alumnos []Alumno
	err = json.Unmarshal(dataAlumnos, &alumnos)
	if err != nil {
		log.Fatal(err)
	}

	for _, alumno := range alumnos {
		val, err := json.Marshal(alumno)
		if err != nil {
			log.Fatal(err)
		}
		if err := CreateUpdate(boltDB, "alumno", []byte(fmt.Sprintf("%d", alumno.Id_alumno)), val); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Alumnos cargados.")

	dataMaterias, err := ioutil.ReadFile("data/materias.json")
	if err != nil {
		log.Fatal(err)
	}

	var materias []Materia
	err = json.Unmarshal(dataMaterias, &materias)
	if err != nil {
		log.Fatal(err)
	}

	for _, materia := range materias {
		val, err := json.Marshal(materia)
		if err != nil {
			log.Fatal(err)
		}
		if err := CreateUpdate(boltDB, "materia", []byte(fmt.Sprintf("%d", materia.Id_materia)), val); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Materias cargadas.")

	dataComisiones, err := ioutil.ReadFile("data/comisiones.json")
	if err != nil {
		log.Fatal(err)
	}

	var comisiones []Comision
	err = json.Unmarshal(dataComisiones, &comisiones)
	if err != nil {
		log.Fatal(err)
	}

	for _, comision := range comisiones {
		val, err := json.Marshal(comision)
		if err != nil {
			log.Fatal(err)
		}
		if err := CreateUpdate(boltDB, "comision", []byte(fmt.Sprintf("%d", comision.Id_comision)), val); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Comisiones cargadas.")
	
	dataCursadas, err := ioutil.ReadFile("data/cursada.json")
	if err != nil {
		log.Fatal(err)
	}

	var cursadas []Cursada
	err = json.Unmarshal(dataCursadas, &cursadas)
	if err != nil {
		log.Fatal(err)
	}

	for _, cursada := range cursadas {
		/*GENERA UN UUID PORQUE LA TABLA CURSADA NO TIENE UN IDENTIFICADOR
		UNICO PARA LA KEY, ENTONCES INSERTA UN UUID PARA CADA KEY -> VALUE*/
		id := uuid.New().String()

		val, err := json.Marshal(cursada)
		if err != nil {
			log.Fatal(err)
		}

		if err := CreateUpdate(boltDB, "cursada", []byte(id), val); err != nil {
			log.Fatal(err)
		}
	}

	fmt.Println("Cursadas cargadas.")
}

func leerBoltDB() {
	boltDB, err := bolt.Open("nosql.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer boltDB.Close()

	fmt.Println("Alumnos:")
	err = boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("alumno"))
		if b == nil {
			fmt.Println("Bucket 'alumno' no encontrado.")
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var alumno Alumno
			if err := json.Unmarshal(v, &alumno); err != nil {
				return err
			}
			fmt.Printf("ID: %s, Nombre: %s, Apellido: %s\n", k, alumno.Nombre, alumno.Apellido)
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Materias:")
	err = boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("materia"))
		if b == nil {
			fmt.Println("Bucket 'materia' no encontrado.")
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var materia Materia
			if err := json.Unmarshal(v, &materia); err != nil {
				return err
			}
			fmt.Printf("ID: %s, Nombre: %s\n", k, materia.Nombre)
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Comisiones:")
	err = boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("comision"))
		if b == nil {
			fmt.Println("Bucket 'comision' no encontrado.")
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var comision Comision
			if err := json.Unmarshal(v, &comision); err != nil {
				return err
			}
			fmt.Printf("ID: %s, Materia ID: %d, Cupo: %d\n", k, comision.Id_materia, comision.Cupo)
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("Cursada:")
	err = boltDB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("cursada"))
		if b == nil {
			fmt.Println("Bucket 'cursada' no encontrado.")
			return nil
		}

		return b.ForEach(func(k, v []byte) error {
			var cursada Cursada
			if err := json.Unmarshal(v, &cursada); err != nil {
				return err
			}
			fmt.Printf("ID Materia: %d, ID Alumno: %d, ID Comision: %d, F Inscripcion: %s, Nota: %d, Estado: %s\n",
				cursada.Id_materia, cursada.Id_alumno, cursada.Id_comision, cursada.F_inscripcion, cursada.Nota, cursada.Estado)
			return nil
		})
	})
	if err != nil {
		log.Fatal(err)
	}
}

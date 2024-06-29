# Trabajo Práctico - Base De Datos 1 - Primer Semestre 2024

**Autores:**  
Joaquin Garcia <joaquin.augusto@hotmail.com>  
Rodrigo Montoro <rodrigo.montoro@hotmail.com>  
Matías Moralez <matiasmoralezc@gmail.com>  
Luciano Rodriguez <luciano.rodriguez2201@mail.com>  

**Docentes:** Hernán Rondelli y Ximena Ebertz  

## Índice de Contenidos

1. [Introducción](#introducción)
2. [Descripción](#descripción)
3. [Implementación](#implementación)
   - [Funciones](#funciones)
     - [inscripcion_materia](#inscripcion_materia)
     - [apertura_inscripcion](#apertura_inscripcion)
     - [baja_de_inscripcion](#baja_de_inscripcion)
     - [cierre_de_inscripcion](#cierre_de_inscripcion)
     - [aplicacion_cupos](#aplicacion_cupos)
     - [ingreso_nota](#ingreso_nota)
     - [cierre_cursada](#ingreso_nota)
   - [Triggers](#triggers)
     - [email_alta_inscripcion_trg](#email_alta_inscripcion_trg)
     - [email_alta_baja_trg](#email_baja_inscripcion_trg)
     - [email_aplicacion_cupos_trg](#email_aplicacion_cupos_trg)
     - [email_inscripcion_lista_espera_trg](#email_inscripcion_lista_espera_trg)
     - [email_cierre_cursada_trg](#email_cierre_cursada_trg)

## Introducción

Este trabajo práctico consiste en modelar una base de datos que simula la inscripción de estudiantes a materias y el registro de sus notas. Para ello se utilizan diferentes funciones y/o triggers, los cuales representan las acciones necesarias para cumplir dichos objetivos. Todo el código debe ser capaz de ejecutarse desde una aplicación CLI escrita en GO y, por último, se solicita guardar determinada información en una base de datos NoSQL basada en Json (utilizando BoltDB), para poder comparar el modelo relacional con uno no relacional NOSQL.

## Descripción

Comenzamos con la creación de un archivo de prueba (que luego fue descartado) para tener una idea de cómo abordar la creación de tablas en lenguaje Go. A partir de ahí planificamos y separamos las funciones necesarias para darle formato a las tablas, otorgarles sus primary y foreign keys, poder borrarlas y cargar los datos desde el archivo Json.

En esta primera parte del TP encontramos como dificultad el reto de aprender el lenguaje Go y adaptarlo al uso de psql, ya que fue un lenguaje nuevo para nosotros. Una vez sentadas las bases, comenzamos con las implementaciones de stored procedures (SP) y triggers que mencionaremos más adelante.

Al momento de realizar los triggers para el envío de emails, surgió una duda que el enunciado no contempla: cada vez que se realice un envío de email, ¿el estado debe quedar pendiente? Ya que si bien se simula el envío de la notificación, no queda claro si debemos interpretar que el correo es enviado. En base a esto también se define el valor de la fecha de envío, por lo que decidimos dejarla en null, ya que no se realiza ningún envío mientras el sistema se ejecute.

## Implementación

A continuación detallamos la implementación de los SP y los triggers. En esta parte del TP encontramos como dificultad que fue necesario "hardcodear" en la tabla dentro de los SP para verificar su correcto funcionamiento dentro del main.

### Funciones

#### inscripcion_materia

Se deberá incluir la lógica que reciba un id de alumno, un id de materia y un id de comisión, y que devuelva true si se logra ingresar la inscripción, o false si se rechaza. El procedimiento deberá validar los siguientes elementos antes de confirmar el alta:

- Que exista un período en estado de inscripción. En caso de que no se cumpla, se debe cargar un error con el mensaje `período de inscripción cerrado`.
- Que el id del alumno exista. En caso de que no se cumpla, se debe cargar un error con el mensaje `id de alumno no válido`.
- Que el id de la materia exista. En caso de que no se cumpla, se debe cargar un error con el mensaje `id de materia no válido`.
- Que el id de comisión exista para la materia. En caso de que no se cumpla, se debe cargar un error con el mensaje `id de comisión no válido para la materia`.
- Que el alumno no esté inscripto previamente en la materia (en cualquiera de sus comisiones). En caso de que no se cumpla, se debe cargar un error con el mensaje `alumno ya inscripto en la materia`.
- Que el alumno tenga en su historia académica todas las materias correlativas en estado regular o aprobada. En caso de que no se cumpla, se debe cargar un error con el mensaje `alumno no cumple requisitos de correlatividad`.

Si se aprueba la solicitud de inscripción, se deberá insertar una fila en la tabla cursada con los datos del alumno, la materia, la comisión y la fecha y hora de inscripción, dejando su estado como ingresado.

#### apertura_inscripcion

Se deberá proveer la lógica que reciba un año y un número de semestre, y que devuelva true si se logra abrir la inscripción para el período o false en caso contrario. El procedimiento deberá validar los siguientes elementos antes de confirmar la apertura:

- Que el año sea mayor o igual al año actual. En caso de que no se cumpla, se debe cargar un error con el mensaje `no se permiten inscripciones para un período anterior`.
- Que el número de semestre sea 1 ó 2. En caso de que no se cumpla, se debe cargar un error con el mensaje `número de semestre no válido`.
- Si el año y semestre solicitado ya existe en la tabla periodo, que su estado sea `cierre inscrip`. En caso de que no se cumpla, se debe cargar un error con el mensaje `no es posible reabrir la inscripción del período, estado actual:[estado]`, reemplazando en el mensaje [estado] por el valor correspondiente que generó el error.
- Que no exista otro período (diferente al solicitado) en estado de inscripción o cierre inscrip. En caso de que no se cumpla, se debe cargar un error con el mensaje `no es posible abrir otro período de inscripción, período actual:[semestre]`, reemplazando en el mensaje [semestre] por el valor correspondiente que generó el error.

Si las validaciones pasan correctamente, se deberá insertar o actualizar la fila correspondiente en la tabla periodo, dejando su estado en inscripción.

#### baja_de_inscripcion

Se deberá proveer la lógica que permita anular la inscripción de un alumno. El procedimiento debe recibir un id de alumno y un id de materia, y retornar true si se logra dar de baja la inscripción o false en caso contrario. El procedimiento deberá validar los siguientes elementos antes de confirmar la baja:

- Que exista un período en estado de inscripción o de cursada. En caso de que no se cumpla, se debe cargar un error con el mensaje `no se permiten bajas en este período`.
- Que el id del alumno exista. En caso de que no se cumpla, se debe cargar un error con el mensaje `id de alumno no válido`.
- Que el id de la materia exista. En caso de que no se cumpla, se debe cargar un error con el mensaje `id de materia no válido`.
- Que el alumno esté inscripto en la materia (en cualquiera de sus comisiones). En caso de que no se cumpla, se debe cargar un error con el mensaje `alumno no inscripto en la materia`.

Si las validaciones pasan correctamente, se deberá actualizar la fila correspondiente de la tabla cursada con el estado `dado de baja`.

En caso de que el período se encuentre en estado de cursada, deberá además actualizarse el estado de un alumno de la misma comisión que se encuentre en espera (si existe alguno), cambiándolo por `aceptado`. Elegir para esto al alumno que tenga la menor fecha de inscripción.

#### cierre_de_inscripcion

Se deberá proveer la lógica que reciba un año y un número de semestre, y que devuelva true si se logra cerrar la inscripción para el período o false en caso contrario. El procedimiento deberá validar los siguientes elementos antes de confirmar el cierre:

- Que el año y semestre solicitado exista en la tabla periodo, y que su estado sea inscripción. En caso de que no se cumpla, se debe cargar un error con el mensaje `el semestre no se encuentra en período de inscripción`.

Si las validaciones pasan correctamente, se deberá actualizar la fila correspondiente en la tabla periodo, dejando su estado en `cierre inscrip`.

#### aplicacion_cupos

Se deberá proveer la lógica que reciba un año y un número de semestre, y que retorne true si se logra aplicar los cupos de cursada a cada comisión o false en caso contrario. El procedimiento deberá realizar las siguientes validaciones antes de aplicar los cupos:

- Que el año y semestre solicitado exista en la tabla periodo, y que su estado sea `cierre inscrip`. En caso de que no se cumpla, se debe cargar un error con el mensaje `el semestre no se encuentra en un período válido para aplicar cupos`.

Si las validaciones pasan correctamente, se deberá asegurar que se ejecuten las siguientes acciones de forma completa (en caso de que se produzca algún error o inconveniente, las acciones que se hayan realizado deberán deshacerse):

- Por cada comisión de materia que tenga alumnos inscriptos, se actualizarán con el estado `aceptado` en la tabla cursada, como máximo la cantidad de alumnos que esté definida en el cupo para la comisión. Al resto de los alumnos de esa comisión, que excedan el cupo, se les actualizará con el estado `en espera`. La prioridad estará dada por el orden de inscripción. A los alumnos que estuvieran dados de baja no se les deberá cambiar su estado, ni se les deberá contar para la aplicación del cupo.
- Luego de que se hayan aplicado los cupos a todas las comisiones, se actualizará la fila de la tabla periodo que tenga el estado `cierre inscrip`, cambiándolo al estado `cursada`.

#### ingreso_nota

Se deberá incluir la lógica que reciba un id de alumno, un id de materia, un id de comisión y una nota, y que devuelva true si se logra ingresar la nota, o false si se rechaza. El procedimiento deberá validar los siguientes elementos antes de confirmar la grabación de la nota:

- Que exista un período en estado de cursada. En caso de que no se cumpla, se debe cargar un error con el mensaje `período de cursada cerrado`.
- Que el id del alumno exista. En caso de que no se cumpla, se debe cargar un error con el mensaje `id de alumno no válido`.
- Que el id de la materia exista. En caso de que no se cumpla, se debe cargar un error con el mensaje `id de materia no válido`.
- Que el id de comisión exista para la materia. En caso de que no se cumpla, se debe cargar un error con el mensaje `id de comisión no válido para la materia`.
- Que el alumno esté inscripto en la comisión de la materia, en estado `aceptado`. En caso de que no se cumpla, se debe cargar un error con el mensaje `alumno no cursa en la comisión`.
- Que la nota se encuentre en el rango de 0 (cero) a 10 (diez), ambos valores inclusive. En caso de que no se cumpla, se debe cargar un error con el mensaje `nota no válida:[nota]`, reemplazando en el mensaje [nota] por el valor que se recibió como parámetro.

Si las validaciones pasan correctamente, se deberá actualizar la nota del alumno en la fila correspondiente de la tabla cursada.

#### cierre_cursada

Se deberá incluir la lógica que reciba un id de materia y un id de comisión, y que devuelva true si se logra completar el cierre de cursada de la comisión, o false si se rechaza. El procedimiento deberá validar los siguientes elementos antes de confirmar el cierre de la cursada:

- Que exista un período en estado de cursada. Si no se cumple, se carga un error con el mensaje `período de cursada cerrado`.
- Que el id de la materia exista. Si no se cumple, se carga un error con el mensaje `id de materia no válido`.
- Que el id de comisión exista para la materia. Si no se cumple, se carga un error con el mensaje `id de comisión no válido para la materia`.
- Que la comisión tenga al menos un alumno inscripto en la tabla `cursada`, sin importar su estado de inscripción. Si no se cumple, se carga un error con el mensaje `comisión sin alumnos inscriptos`.
- Que todos los alumnos de la comisión en estado 'aceptado' tengan informada su nota de cursada. Si no se cumple, se carga un error con el mensaje `la carga de notas no está completa`.

Si todas las validaciones son exitosas, se ejecutan las siguientes acciones de forma completa:
- Se inserta una fila en la tabla `historia_academica` por cada alumno de la comisión que se encuentre en estado 'aceptado'. El semestre corresponde al período con estado de cursada, y el estado de la historia académica dependerá de la nota de regularidad (0: ausente, 1-3: reprobada, 4-6: regular, 7-10: aprobada). La nota final se graba solo si el estado es 'aprobada', siendo igual a la nota de regularidad en ese caso.
- Después de insertar en la historia académica todos los alumnos en estado 'aceptado', se eliminan de la tabla `cursada` todos los registros de esa comisión, independientemente de su estado.

### Triggers

#### email_alta_inscripcion_trg

Este trigger se dispara cada vez que se realiza una actualización en la tabla `cursada` y se fija si el estado de la inscripción toma el valor 'ingresada'. Si esto se cumple, ejecuta la función `email_alta_inscipcion()` la cual genera una nueva fila en `envio_mail` indicando el asunto, la fecha de creación de la fila, el email del alumno, la fecha en la que se envió el correo, el estado y el cuerpo del email con los datos de la materia, la comisión y del alumno, notificando que la inscripción se ha registrado en el sistema.

#### email_baja_inscripcion_trg

Este trigger se dispara cada vez que se realiza una actualización en la tabla `cursada` y se fija si el estado anterior a la actualización de la inscripción era 'ingresada' y si el nuevo estado toma los valores: 'aceptado' o 'en espera'. Si la condición se cumple, el trigger ejecuta la función `email_aplicacion_cupos()` la cual inserta una nueva fila en `envio_mail` indicando el asunto, la fecha de creación de la fila, el email del alumno, la fecha en la que se envió el correo, el estado y el cuerpo del email con los datos de la materia, la comisión y del alumno, además informando el estado en el que se encuentra su inscripción.

#### email_aplicacion_cupos_trg

Este trigger se dispara cada vez que se realiza una actualización en la tabla `cursada` y verifica si el estado anterior a la actualización de la inscripción era 'ingresada' y si el nuevo estado toma los valores: 'aceptado' o 'en espera'. Si la condición se cumple, el trigger ejecuta la función `email_aplicacion_cupos()` la cual inserta una nueva fila en `envio_mail` indicando el asunto, la fecha de creación de la fila, el email del alumno, la fecha en la que se envió el correo, el estado y el cuerpo del email con los datos de la materia, la comisión y del alumno, además informando el estado en el que se encuentra su inscripción.

#### email_inscripcion_lista_espera_trg

Este trigger se dispara cada vez que se realiza una actualización en la tabla `cursada` y se verifica si el estado anterior a la actualización de la inscripción era 'en espera' y si el nuevo estado toma el valor 'aceptada'. Si la condición se cumple, el trigger ejecuta la función `email_inscripcion_lista_espera()` la cual inserta una nueva fila en `envio_mail` indicando el asunto, la fecha de creación de la fila, el email del alumno, la fecha en la que se envió el correo, el estado y el cuerpo del email donde se informa el cambio de estado de la inscripción, mostrando además los datos de la materia, la comisión y del alumno.

#### email_cierre_cursada_trg

Este trigger se dispara cada vez que se realiza una actualización en la tabla `cursada` y se fija si el estado tiene el valor 'aceptado'. Si se cumple la condición, el trigger ejecuta la función `email_cierre_cursada()` la cual hace un insert en `envio_mail` indicando el asunto, la fecha de creación de la fila, el email del alumno, la fecha en la que se envió el correo, el estado y el cuerpo del email donde se informa el cambio de estado de la cursada junto al estado académico del alumno, su nota regular y nota final en caso de tener.

### Conclusiones

En este trabajo práctico, asumimos el desafío de modelar una base de datos para la inscripción de alumnos en materias y el registro de sus notas a través de un desarrollo colaborativo. Para lograrlo, tuvimos que aprender y adaptar el lenguaje Go al uso de PostgreSQL, lo cual representó un reto y una valiosa oportunidad de aprendizaje.

Durante la implementación, trabajamos en varias funciones y triggers que emulan las operaciones necesarias dentro del sistema de inscripciones. Ajustamos manualmente diversas partes del código para asegurar su correcto funcionamiento y para que las diferentes acciones y actualizaciones en la base de datos se realizaran de manera adecuada.

En resumen, este proyecto nos permitió aprender a combinar el uso de bases de datos relacionales y no relacionales, mejorar nuestras habilidades en Go y PostgreSQL, y desarrollar soluciones prácticas para la gestión de inscripciones y notas en un entorno académico. Aunque presentó sus desafíos, fue una experiencia enriquecedora que facilitó nuestro crecimiento en el ámbito del desarrollo y nos proporcionó un mejor entendimiento de la integración de diversas tecnologías.

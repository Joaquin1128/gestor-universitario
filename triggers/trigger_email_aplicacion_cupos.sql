create or replace function email_aplicacion_cupos() returns trigger as $$
declare
	v_materia_nombre text;
	v_comision_numero text;
	v_alumno_nombre text;
	v_alumno_apellido text;
	v_email_alumno text;
	v_estado_inscripcion char(12);
begin
	select nombre into v_materia_nombre from materia where id_materia = new.id_materia;
	select id_comision into v_comision_numero from comision where id_materia = new.id_materia and id_comision = new.id_comision;
	select nombre, apellido, email into v_alumno_nombre, v_alumno_apellido, v_email_alumno from alumno where id_alumno = new.id_alumno;
		
	if new.estado = 'aceptada' then
		v_estado_inscripcion := 'aceptada';
	else
		v_estado_inscripcion := 'en espera';
	end if;
		
	insert into envio_email
	values (nextval('envio_email_id_seq'), current_timestamp, v_email_alumno, v_estado_inscripcion,
	'Hola ' || v_alumno_nombre || ' ' || v_alumno_apellido || ', tu inscripcion a la materia ' || v_materia_nombre || ', comision ' || v_comision_numero ||' se encuentra: ' || v_estado_inscripcion || '.',
	null, 'pendiente'
	);
		
	return new;
end;
$$ language plpgsql;

create trigger email_aplicacion_cupos_trg
after update on cursada
for each row
when (old.estado = 'ingresada' and (new.estado = 'aceptada' or new.estado = 'en espera'))
execute	function email_aplicacion_cupos();

create or replace function email_inscripcion_lista_espera() returns trigger as $$
declare
	v_materia_nombre text;
	v_comision_numero text;
	v_alumno_nombre text;
	v_alumno_apellido text;
	v_email_alumno text;
begin
	select nombre into v_materia_nombre from materia where id_materia = new.id_materia;
	select id_comision into v_comision_numero from comision where id_materia = new.id_materia and id_comision = new.id_comision;
	select nombre, apellido, email into v_alumno_nombre, v_alumno_apellido, v_email_alumno from alumno where id_alumno = new.id_alumno;
	
	insert into envio_email
	values (nextval('envio_email_id_seq'), current_timestamp, v_email_alumno, 'Inscripcion aceptada',
	'Hola ' || v_alumno_nombre || ' ' || v_alumno_apellido || ', tu inscripcion a la materia ' || v_materia_nombre || ', comision ' || v_comision_numero ||' ha sido finalmente aceptada',
	null, 'pendiente'
	);
	
	return new;
end;
$$ language plpgsql;
	
create trigger email_inscripcion_lista_espera_trg
after update on cursada
for each row
when (old.estado = 'en espera' and new.estado = 'aceptada')
execute function email_inscripcion_lista_espera();

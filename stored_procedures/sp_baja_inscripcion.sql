create function baja_de_inscripcion(id_alumno_buscado integer, id_materia_buscada integer, out p_result boolean, out p_error_message text) as $$
declare
	resultado_periodo periodo%rowtype;
	resultado_alumno alumno%rowtype;
	resultado_materia materia%rowtype;
	resultado_comision comision%rowtype;
	resultado_cursada cursada%rowtype;

	alumno_enespera record;

begin
	select * into resultado_periodo from periodo where estado = 'inscripcion' or estado = 'cursada';

	if not found then
		insert into error values(nextval('error_id_seq'), 'baja inscrip', null, id_alumno_buscado, id_materia_buscada, null, current_timestamp, 'No se permiten bajas en este periodo');
		p_error_message := 'no se permiten bajas en este periodo';
		p_result := false;
		return;
	end if;

	select * into resultado_alumno from alumno where id_alumno = id_alumno_buscado;

	if not found then
		insert into error values(nextval('error_id_seq'), 'baja inscrip', resultado_periodo.semestre, id_alumno_buscado, id_materia_buscada, null, current_timestamp, 'Id de alumno no válido');
		p_error_message := 'id de alumno no válido';
		p_result := false;
		return;
	end if;

	select * into resultado_materia from materia where id_materia = id_materia_buscada;

	if not found then
		insert into error values(nextval('error_id_seq'), 'baja inscrip', resultado_periodo.semestre, id_alumno_buscado, id_materia_buscada, null, current_timestamp, 'Id de materia no válido');
		p_error_message := 'id de materia no válido';
		p_result := false;
		return;
	end if;

	select * into resultado_cursada from cursada where id_alumno = id_alumno_buscado and id_materia = id_materia_buscada and estado = 'aceptade';

	if not found then
		insert into error values(nextval('error_id_seq'), 'baja inscrip', resultado_periodo.semestre, id_alumno_buscado, id_materia_buscada, null, current_timestamp, 'alumno no inscripte en la materia');
		p_error_message := 'alumno no inscripto en la materia';
		p_result := false;
		return;
	end if;

	update cursada set estado = 'dada de baja' where cursada.id_alumno = id_alumno_buscado and cursada.id_materia = id_materia_buscada;
	
	if resultado_periodo.estado = 'cursada' then
		select * into alumno_enespera from cursada 
		where id_materia = id_materia_buscada and id_comision = resultado_cursada.id_comision and estado = 'en espera' 
		order by f_inscripcion asc limit 1;
		
		update cursada set estado = 'aceptada' 
		where id_alumno = alumno_enespera.id_alumno and id_materia = id_materia_buscada and id_comision = resultado_cursada.id_comision;
	end if;
end;
$$ language plpgsql;

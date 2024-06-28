create function aplicacion_cupos(anio_ingresado int, nro_semestre_ingresado int, out p_result boolean, out p_error_message text) as $$
declare
	periodo_encontrado periodo%rowtype;
	cupo_materia int;
	id_materia_buscada int := 1;
	id_comision_buscada int;
	alumno_inscripte cursada%rowtype;
	comision_materia comision%rowtype;
	materia comision%rowtype;
	semestre_buscado varchar(6);
begin
	semestre_buscado := to_char(anio_ingresado, 'FM999999') || '-' || to_char(nro_semestre_ingresado, 'FM999999');

	select * into periodo_encontrado from periodo where semestre = semestre_buscado and estado = 'cierre inscrip';
	
	if not found then
		insert into error values(nextval('error_id_seq'), 'aplicacion cupo', semestre_buscado, null, null, null, current_timestamp, 'El semestre no se encuentra en un período válido para aplicar cupos');
		p_result := false;
		p_error_message := 'El semestre no se encuentra en un período válido para aplicar cupos';
		return;
	end if;
	
	loop
		perform 1 from comision where id_materia = id_materia_buscada;
		exit when not found;
	
		for comision_materia in select * from comision where id_materia = id_materia_buscada loop
			id_comision_buscada := comision_materia.id_comision;
			
			select cupo into cupo_materia from comision where id_materia = id_materia_buscada and id_comision = comision_materia.id_comision;
			
			for alumno_inscripte in (select * from cursada where id_materia = id_materia_buscada and id_comision = id_comision_buscada and estado = 'ingresada' order by f_inscripcion asc limit cupo_materia) loop
				update cursada set estado = 'aceptada' 
				where id_alumno = alumno_inscripte.id_alumno 
				and id_materia = alumno_inscripte.id_materia 
				and id_comision = alumno_inscripte.id_comision;
			end loop;
		
		update cursada set estado = 'en espera'
		where id_materia = id_materia_buscada and id_comision = comision_materia.id_comision and estado = 'ingresade';
		
		end loop;
		
		id_materia_buscada := id_materia_buscada + 1;
	end loop;
	
	update periodo set estado = 'cursada'
	where estado = 'cierre inscrip';
	
	p_result := true;
end;
$$ language plpgsql;

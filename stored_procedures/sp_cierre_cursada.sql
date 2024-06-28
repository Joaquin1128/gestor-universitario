create or replace function cierre_cursada(id_materia_buscada int, id_comision_buscada int, out p_result boolean, out p_error_message text) as $$
declare
	v_count int;
	v_semestre_actual text;
	v_nota_regular int;
	v_nota_final int;
begin	
	p_error_message := '';
	
	select semestre into v_semestre_actual from periodo where estado = 'cursada' limit 1;

	if v_semestre_actual is null then
		insert into error values(nextval('error_id_seq'), 'cierre cursada', v_semestre_actual, null, id_materia_buscada, id_comision_buscada, current_timestamp, 'Periodo de cursada cerrado');
		p_result := false;
		p_error_message := 'periodo de cursada cerrado';
		return;
	end if;
	
	if not exists (select 1 from materia where id_materia = id_materia_buscada) then
		insert into error values(nextval('error_id_seq'), 'cierre cursada', v_semestre_actual, null, id_materia_buscada, id_comision_buscada, current_timestamp, 'Id de materia no valido');
        p_result := false;
        p_error_message := 'id de materia no valido';
		return;
	end if;
  
	if not exists (
		select 1 from comision
		where id_materia = id_materia_buscada and
		id_comision = id_comision_buscada
		) then
			insert into error values(nextval('error_id_seq'), 'cierre cursada', v_semestre_actual, null, id_materia_buscada, id_comision_buscada, current_timestamp, 'Id de comision no válido para la materia');
			p_result := false;
			p_error_message := 'id de comision no valido para la materia';
			return;
	end if;
	
	select count(*) into v_count from cursada where id_comision = id_comision_buscada and
	id_materia = id_materia_buscada;
	
	if v_count = 0 then
		insert into error values(nextval('error_id_seq'), 'cierre cursada', v_semestre_actual, null, id_materia_buscada, id_comision_buscada, current_timestamp, 'Comision sin alumnos inscriptos');
		p_result := false;
		p_error_message := 'comision sin alumnos inscriptos';
		return;
	end if;
	
	select count(*) into v_count from cursada where id_comision = id_comision_buscada and
	estado = 'aceptada' and
	nota is null;
	
	if v_count != 0 then
		insert into error values(nextval('error_id_seq'), 'cierre cursada', v_semestre_actual, null, id_materia_buscada, id_comision_buscada, current_timestamp, 'La carga de notas no esta completa');
		p_result := false;
		p_error_message := 'la carga de notas no esta completa';
		return;
	end if;
	
	insert into historia_academica (id_alumno, semestre, id_materia, id_comision, estado, nota_regular, nota_final)
	select c.id_alumno, v_semestre_actual, c.id_materia, c.id_comision,
		case
			when c.nota is null then 'ausente'
			when c.nota between 1 and 3 then 'reprobada'
			when c.nota between 4 and 6 then 'regular'
			when c.nota between 7 and 10 then 'aprobada'
		end as estado,
		c.nota as nota_regular,
		case 
			when c.nota between 7 and 10 then c.nota
			else null
		end as nota_final
	from cursada c
	where c.id_comision = id_comision_buscada and
	c.id_materia = id_materia_buscada and
	c.estado = 'aceptada';
	
	delete from cursada
	where id_comision = id_comision_buscada and
	id_materia = id_materia_buscada;
	
	p_result := true;
	
end;
$$ language plpgsql;

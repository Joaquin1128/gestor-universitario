create or replace function ingreso_nota(id_alumno_buscado int, id_materia_buscada int, id_comision_buscada int, nota_ingresada int, out p_result boolean, out p_error_message text) as $$
declare
	v_periodo periodo%rowtype;
begin
	p_error_message := '';

	select * into v_periodo from periodo where estado = 'cursada';

	if not found then
		insert into error values(nextval('error_id_seq'), 'ingreso nota', v_periodo.semestre, id_alumno_buscado, id_materia_buscada, id_comision_buscada, current_timestamp, 'La carga de notas no esta completa');
		p_result := false;
		p_error_message := 'la carga de notas no esta completa';
		return;
	end if;
	
	if not exists (select 1 from alumno where id_alumno = id_alumno_buscado) then
		insert into error values(nextval('error_id_seq'), 'ingreso nota', v_periodo.semestre, id_alumno_buscado, id_materia_buscada, id_comision_buscada, current_timestamp, 'Id de alumno no valido');
        p_result := false;
        p_error_message := 'id de alumno no valido';
		return;
	end if;
	
	if not exists (select 1 from materia where id_materia = id_materia_buscada) then
		insert into error values(nextval('error_id_seq'), 'ingreso nota', v_periodo.semestre, id_alumno_buscado, id_materia_buscada, id_comision_buscada, current_timestamp, 'Id de materia no valido');
        p_result := false;
        p_error_message := 'id de materia no valido';
		return;
	end if;
  
	if not exists (
		select 1 from comision
		where id_materia = id_materia_buscada and
		id_comision = id_comision_buscada
		) then
			insert into error values(nextval('error_id_seq'), 'ingreso nota', v_periodo.semestre, id_alumno_buscado, id_materia_buscada, id_comision_buscada, current_timestamp, 'Id de comision no valido para la materia');
			p_result := false;
			p_error_message := 'id de comision no valido para la materia';
			return;
	end if;
	
	if not exists (
		select 1 from cursada
		where id_alumno = id_alumno_buscado and
		id_materia = id_materia_buscada and
		id_comision = id_comision_buscada and
		estado = 'aceptade'
		) then
			insert into error values(nextval('error_id_seq'), 'ingreso nota', v_periodo.semestre, id_alumno_buscado, id_materia_buscada, id_comision_buscada, current_timestamp, 'alumno no cursa en la comision');
			p_result := false;
			p_error_message := 'alumno no cursa en la comision';
			return;
	end if;
	
	if nota_ingresada < 0 or nota_ingresada > 10 then
		insert into error values(nextval('error_id_seq'), 'ingreso nota', v_periodo.semestre, id_alumno_buscado, id_materia_buscada, id_comision_buscada, current_timestamp, 'Nota no válida');
		p_result := false;
		p_error_message := 'nota no valida: ' || nota;
		return;
	end if; 
	
	update cursada 
	set nota = nota_ingresada
	where id_alumno = id_alumno_buscado 
	and id_materia = id_materia_buscada 
	and id_comision = id_comision_buscada;

	p_result := true;
	
end;
$$ language plpgsql;

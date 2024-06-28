create or replace function apertura_inscripcion(anio_ingresado int, nro_semestre_ingresado int, out p_result boolean, out p_error_message text) as $$
declare
	v_estado_actual varchar(6);
	v_anio_actual int;
	v_count int;
	v_nro_semestre char(1);
	v_semestre char(6);
begin
	p_error_message := '';

	v_anio_actual := anio_ingresado;
	v_nro_semestre := to_char(nro_semestre_ingresado, 'FM999999');
	v_semestre := to_char(anio_ingresado, 'FM999999') || '-' || v_nro_semestre;

	if v_nro_semestre not in ('1', '2') then
		insert into error values(nextval('error_id_seq'), 'apertura', v_semestre, null, null, null, current_timestamp, 'Número de semestre no válido');
		p_result := false;
		p_error_message := 'número de semestre no válido';
		return;
	end if;

	if v_anio_actual < extract(year from current_date) then
		insert into error values(nextval('error_id_seq'), 'apertura', v_semestre, null, null, null, current_timestamp, 'No se permiten inscripciones para un período anterior');
		p_result := false;
		p_error_message := 'no se permiten inscripciones para un período anterior';
		return;
	end if;

	select estado into v_estado_actual from periodo where semestre = v_semestre;

	if v_estado_actual is not null and v_estado_actual != 'cerrado' then
		insert into error values(nextval('error_id_seq'), 'apertura', v_semestre, null, null, null, current_timestamp, 'No es posible reabrir la inscripción del período');
		p_result := false;
		p_error_message := format('no es posible reabrir la inscripción del período, estado actual: %s', v_estado_actual);
		return;
	end if;

	select count(*) into v_count from periodo where estado in ('inscripcion', 'cierre inscrip') and semestre != v_semestre;

	if v_count > 0 then
		insert into error values(nextval('error_id_seq'), 'apertura', v_semestre, null, null, null, current_timestamp, 'No es posible abrir otro período de inscripción, ya existe otro período en estado inscripción o cierre inscripción');
		p_result := false;
		p_error_message := 'no es posible abrir otro período de inscripción, ya existe otro período en estado inscripción o cierre inscripción';
		return;
	end if;

	insert into periodo (semestre, estado) values (v_semestre, 'inscripcion')
	on conflict (semestre) do update set estado = excluded.estado;

	p_result := true;
end;
$$ language plpgsql;

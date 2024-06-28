create function cierre_de_inscripcion(anio_ingresado int, nro_semestre_ingresado int, out p_result boolean, out p_error_message text) as $$
declare
	resultado_periodo periodo%rowtype;
	semestre_buscado varchar(6);
begin
	semestre_buscado := to_char(anio_ingresado, 'FM999999') || '-' || to_char(nro_semestre_ingresado, 'FM999999');

	select * into resultado_periodo from periodo where semestre = semestre_buscado and estado = 'inscripcion' ;

	if not found then
		insert into error values(nextval('error_id_seq'), 'aplicacion cupo', semestre_buscado, null, null, null, current_timestamp, 'El semestre no existe en periodo de inscripcion');
		p_result := false;
		p_error_message := 'El semestre no se encuentra en un período válido para aplicar cupos';
		return;
	end if;

	update periodo set estado = 'cierre inscrip' where semestre = semestre_buscado;
	
	p_result := true;
end;
$$ language plpgsql;

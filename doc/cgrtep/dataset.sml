dataset main
    table fm
    table f
    table locations
    join fm <-> f        
        on f.id = fm.family_id
    join f <-> locations
        on f.location_id = locations.id 
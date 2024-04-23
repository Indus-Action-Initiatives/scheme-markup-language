scheme sewing_machine
    label Chief Minister Sewing Machine Assistance Scheme
    description """Female construction worker receives sewing machine"""    
    criteria bocw_card
        column has_bocw_card
        table fm
        operator equals
        value true
    criteria age       
        column dob
        table fm
        operator age_between
        value [18, 50]
        granularity year
    criteria gender
        column gender
        table fm
        operator equals
        value female
    
    evaluation bocw_card && age && gender
        
scheme uow_cycle
    label Chief Minister Unorganized Workers Cycle Assistance Scheme
    description """One cycle (per beneficiary) will be payable."""    
    criteria uow_card
        column has_uow_card
        table fm
        operator equals
        value true
    criteria age
        column dob
        table fm
        operator age_between
        value [18, 40]
        granularity year
    criteria gender
        column gender
        table fm
        operator equals
        value female
    
    evaluation bocw_card && age && gender
        
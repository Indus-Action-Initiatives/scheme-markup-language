scheme family_pension
    label Family Pension
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria health_status
        table fm
        column health_status
        operator equals
        value Deceased
    criteria spouse_alive
        table fm
        column spouse_alive
        operator equals
        value True
    criteria spouse_occupation
        table fm
        column spouse_occupation
        operator equals
        value Pensioner
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte // generalising age to be difference between the column and today
        granularity year
        value 1

    evaluation occupation && health_status && spouse_alive && spouse_occupation && card_registration
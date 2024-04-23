scheme old_age_pension
    label Old Age Pension Benefit
    criteria age
        table fm
        column age
        operator gt
        value 60
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria health_status
        table fm
        column health_status
        operator ne
        value Deceased
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte // generalising age to be difference between the column and today
        granularity year
        value 3

    evaluation age && occupation && health_status && card_registration
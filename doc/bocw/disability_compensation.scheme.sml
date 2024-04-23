scheme disability_compensation
    label Compensation for disability
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria health_status
        table fm
        column fm
        operator equals
        value Permanent disability as per disability certificate
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte // generalising age to be difference between the column and today
        granularity year
        value 0

    evaluation occupation && health_status && card_registration
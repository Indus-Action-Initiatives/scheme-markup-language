scheme medical_assistance
    label Medical Assistance
    description """Financial Assistance for marriage of self and for children (upto 2 children). (Rule – 282), the building workers having continuous membership of 03 years shall be eligible. The details are as under:-
Marriage of female registered member – Rs.51,000/-
Marriage of male registered member    - Rs.35,000/-
Marriage of daughter of registered members – Rs.51,000/-
Marriage of son of registered members – Rs.35,000/
and wives of male members (upto 2 children). (Rule – 271) – from the date of joining membership of the fund."""    
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria health_status
        table fm
        column health_status
        operator IN
        value ['Hospitalised for more than 5 days', 'In plaster at residence']
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte // generalising age to be difference between the column and today
        granularity year
        value 0

    evaluation occupation && health_status && card_registration
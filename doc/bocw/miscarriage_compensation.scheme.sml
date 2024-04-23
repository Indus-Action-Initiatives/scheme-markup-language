scheme miscarriage_compensation
    label Compensation in case of miscarriage
    criteria gender_and_age
        combine OR
            term male_gender_and_age
                combine AND
                    term male_gender
                        table fm
                        column gender
                        operator equals
                        value male
                    term male_age
                        table fm
                        column age
                        operator gte
                        value 21   
            term female_gender_and_age         
                combine AND
                    term female_gender
                        table fm
                        column gender
                        operator equals
                        value female
                    term female_age
                        table fm
                        column age
                        operator gte
                        value 18
    criteria occupation
        table fm
        column occupation
        operator equals
        value Construction Worker
    criteria marital_status
        table fm
        column marital_status
        operator equals
        value Married
    criteria pregnancy_status
        table fm
        column pregnancy_status
        operator IN
        value ['Miscarriage', 'Still born']
    criteria card_registration
        table fm
        column bocw_card_registration_date
        operator age_gte
        granularity year
        value 0

    evaluation age && occupation && marital_status && pregnancy_status && card_registration
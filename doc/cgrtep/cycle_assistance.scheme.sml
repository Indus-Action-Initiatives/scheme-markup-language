scheme cycle_assistance
    label Chief Minister Cycle Assistance Scheme
    description """Construction worker receives cycle"""    
    criteria bocw_card        
        column has_bocw_card
        table fm
        operator equals
        value true
    criteria age_by_gender
        combine OR
            term female
                combine AND
                    term female_age
                        column dob
                        table fm
                        operator age_between
                        value [18, 35]
                        granularity year
                    term female_gender
                        column gender
                        table fm
                        operator equals
                        value female
            term male
                combine AND
                    term male_age
                        column dob
                        table fm
                        operator age_between
                        value [18, 50]
                        granularity year
                    term male_gender
                        column gender
                        table fm
                        operator equals
                        value male
    
    evaluation bocw_card && age_by_gender
        
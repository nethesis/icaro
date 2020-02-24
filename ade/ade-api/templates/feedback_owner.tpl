<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <style>
        h1 {
            color: {{ .TitleColor }};
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        h2 {
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
            margin: 20px;
        }

        h3 {
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        h4 {
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        h5 {
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        img {
            width: 150px;
        }

        .star {
            width: 50px;
            margin-bottom: 25px;
        }

        .title {
            margin: 20px;
        }

        .container {
            background: {{ .ContainerBgColor }};
            margin-top: 20px !important;
            margin-bottom: 20px !important;
            margin-left: auto !important;
            margin-right: auto !important;
            border-radius: .28571429rem;
            padding: 10px;
            width: 70%;
        }

        .divider {
            padding: 5px;
        }
    </style>
    <title>Review</title>
</head>

<body style="margin: 0;text-align: center;font-family: Lato,'Helvetica Neue',Arial,Helvetica,sans-serif; color: {{ .TextColor }} !important;">
    <table style="width: 100%; text-align: center; border-collapse: collapse; background-color: {{ .BgColor }};">
        <tbody>
            <tr>
                <td style="width: 100%;">
                    <div class="container">
                        <div class="title">
                            <h3>Un cliente ha inviato una segnalazione:</h3>
                            <p>
                                <div><em>
                                    {{ .Message }}
                                </em></div>
                            </p>
                            <small>Questa è una mail automatica. Non rispondere</small>
                        </div>
                        <div class="divider"></div>
                        <div>
                            <h3>A customer sent a report:</h3>
                            <p>
                                <div><em>
                                    {{ .Message }}
                                </em></div>
                            </p>
                            <small>This is an automatic email. Do not reply</small>
                        </div>
                    </div>
                </td>
            </tr>
        </tbody>
    </table>
</body>

</html>

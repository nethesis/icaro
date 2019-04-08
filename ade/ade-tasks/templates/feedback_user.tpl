<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <style>
        h1 {
            color: white;
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
            color: white;
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        .title {
            margin: 20px;
        }

        .container {
            background: white;
            margin-left: 10px !important;
            margin-right: 10px !important;
            border-radius: .28571429rem;
            padding: 10px;
        }

        .divider {
            padding: 5px;
        }
    </style>
    <title>Feedback</title>
</head>

<body style="margin: 0;text-align: center;font-family: Lato,'Helvetica Neue',Arial,Helvetica,sans-serif;">
    <table style="width: 100%; border-collapse: collapse; background-color: {{ .BgColor }};">
        <tbody>
            <tr>
                <td style="width: 100%;">
                    <div class="title">
                        <h1>{{ .HotspotName }}</h1>
                        <img style="width: 150px;" src="{{ .Base64Type }},{{ .HotspotLogo }}">
                    </div>
                    <div class="container">
                        {{ .HotspotSurveyBodyText }}
                    </div>
                    <h5>
                        {{ .HotspotDetails }}
                    </h5>
                </td>
            </tr>
        </tbody>
    </table>
</body>

</html>

<!DOCTYPE html>
<html>

<head>
    <meta charset="UTF-8">
    <style>
        body {
            background-color: {{ .BgColor }};
            text-align: center;
            font-family: Lato,
            'Helvetica Neue',
            Arial,
            Helvetica,
            sans-serif;
        }

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
            border: none;
            padding: 0 0;
            font-weight: 700;
            text-transform: none;
        }

        p {
            color: red;
        }

        img {
            width: 150px;
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

<body>
    <div class="title">
        <h1>{{ .HotspotName }}</h1>
        <img src="data:image/png;base64,{{ .HotspotLogo }}">
    </div>
    <div class="container">
        {{ .HotspotFeedbackBodyText}}
        <h5>
            {{ .HotspotDetails }}
        </h5>
    </div>

</body>

</html>
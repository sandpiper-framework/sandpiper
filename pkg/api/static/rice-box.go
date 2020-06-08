package static

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file3 := &embedded.EmbeddedFile{
		Filename:    "css/style.css",
		FileModTime: time.Unix(1591640285, 0),

		Content: string("@import url(https://fonts.googleapis.com/css?family=Roboto:300);\r\n\r\n.register-page {\r\n  width: 360px;\r\n  padding: 8% 0 0;\r\n  margin: auto;\r\n}\r\n.form {\r\n  position: relative;\r\n  z-index: 1;\r\n  background: #FFFFFF;\r\n  max-width: 360px;\r\n  margin: 0 auto 100px;\r\n  padding: 45px;\r\n  text-align: center;\r\n  box-shadow: 0 0 20px 0 rgba(0, 0, 0, 0.2), 0 5px 5px 0 rgba(0, 0, 0, 0.24);\r\n}\r\n.form input, select {\r\n  font-family: \"Roboto\", sans-serif;\r\n  outline: 0;\r\n  background: #f2f2f2;\r\n  width: 100%;\r\n  border: 0;\r\n  margin: 0 0 15px;\r\n  padding: 15px;\r\n  box-sizing: border-box;\r\n  font-size: 14px;\r\n}\r\n.form button {\r\n  font-family: \"Roboto\", sans-serif;\r\n  text-transform: uppercase;\r\n  outline: 0;\r\n  background: rgb(61, 72, 122);\r\n  width: 100%;\r\n  border: 0;\r\n  padding: 15px;\r\n  color: #FFFFFF;\r\n  font-size: 14px;\r\n  -webkit-transition: all 0.3s ease;\r\n  transition: all 0.3s ease;\r\n  cursor: pointer;\r\n}\r\n.form button:hover,.form button:active,.form button:focus {\r\n  background: #43A047;\r\n}\r\n.form .message {\r\n  margin: 15px 0 0;\r\n  color: #b3b3b3;\r\n  font-size: 12px;\r\n}\r\n.form .message a {\r\n  color: #b3b3b3;\r\n  text-decoration: none;\r\n}\r\n.container {\r\n  position: relative;\r\n  z-index: 1;\r\n  max-width: 300px;\r\n  margin: 0 auto;\r\n}\r\n.container:before, .container:after {\r\n  content: \"\";\r\n  display: block;\r\n  clear: both;\r\n}\r\n.container .info {\r\n  margin: 50px auto;\r\n  text-align: center;\r\n}\r\n.container .info h1 {\r\n  margin: 0 0 15px;\r\n  padding: 0;\r\n  font-size: 36px;\r\n  font-weight: 300;\r\n  color: #1a1a1a;\r\n}\r\n.container .info span {\r\n  color: #4d4d4d;\r\n  font-size: 12px;\r\n}\r\n.container .info span a {\r\n  color: #000000;\r\n  text-decoration: none;\r\n}\r\n.container .info span .fa {\r\n  color: #EF3B3A;\r\n}\r\nbody {\r\n  background: rgb(58, 78, 168); /* fallback for old browsers */\r\n  background: -webkit-linear-gradient(right, #060b22, rgb(58, 78, 168));\r\n  background: -moz-linear-gradient(right, #060b22, rgb(58, 78, 168));\r\n  background: -o-linear-gradient(right, #060b22, rgb(58, 78, 168));\r\n  background: linear-gradient(to left, #060b22, rgb(58, 78, 168));\r\n  font-family: \"Roboto\", sans-serif;\r\n  -webkit-font-smoothing: antialiased;\r\n  -moz-osx-font-smoothing: grayscale;      \r\n}\r\n.logo {\r\n  display: block;\r\n  margin-left: auto;\r\n  margin-right: auto;\r\n  width: 75%;\r\n  margin-bottom: 20px;\r\n}"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "images/favicon-16x16.png",
		FileModTime: time.Unix(1591403053, 0),

		Content: string("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x10\x00\x00\x00\x10\b\x06\x00\x00\x00\x1f\xf3\xffa\x00\x00\x01\x02IDAT8O\x8d\x93}q\x02A\fG\x1f\n\x8a\x03Z\x05\x05\x05P\x05\x80\x82\xd6\x01T\x018hqP\x14\x80\x03@A\xc1A\x1d\x14\t\xcc\xebl\x98\xe5z\xb7\\f\xee\x0f\xf2\xf16\xf9%t\xf8o\x8f\xc0\f\x18\x01\xfd\x14\xfe\x01\xf6\xc0\n8\xe6%\x9d\xecG\x17X\x00\xf3\x1ah\xee\xfa\x04\xde\xc3\x11\x00\x8bwًw\x18\u007f]\fL\n\x80T\xdb\x0e;\x00\xfa\xce\xc0\x1b\xf0ZCt\x9c\xb9\x00\xe7\xfc\xae\x14;\u007fn\xcb4^\x95\xf3$\xe0\xab\xf2\xc2\x14ئL\x85S\x13E\xfc\xad\xe9b-\xc0`/\v\xbe$\xc5uٝqG\x114I\xbe\x87\x94\u007f\x8c\x11b5\xb6n\x81_\xc9\x04\xab\xcd,_c\xa9@\xf00݆yW\x8d\xda\x00l{\xd3@?\xb5\x01(\xe8\xb8\x01\xb0j\x03P\xc0\x10\xadv\x8d\xa5\xd9K\xed_\x0f\xa9\x04\xf0\x0e\x14\xafj\xa78\xfb\xd2\b\xee\xfd\xa3\xe9\x84\xc3\xdf\x04\xf0/\xadx\xcf)\xd1\x17\xbd\x15O\xfa\xe6F.\xa0\xfd/\x1eK\xff\t\x13\x00\x00\x00\x00IEND\xaeB`\x82"),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "images/favicon-32x32.png",
		FileModTime: time.Unix(1591403053, 0),

		Content: string("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00 \x00\x00\x00 \b\x06\x00\x00\x00szz\xf4\x00\x00\x02|IDATXG\xb5\xd7K\xc8MQ\x14\xc0\xf1\xdf7 \x8c<\xf2,D\f\x94G\x89B\x14I!$\xa5\x18\b呈\xa4d\x86\x91PJ\x1e\x03\xaf\x91\"\x06\x842\xf1\x1a\x11\x92G\x91\xc7\x00\x13&H^)%\xad\xaf\xf3\xddN\xb7{\xee\xdd\xe7\xbbǪݹݻ\x1e\xff\xbb\xf6\xdek\xadӡ\x9c,\xc4\x12L\xc6\xc0l\xfd\xc5G|\xc03<\xc6m\xbcOqݑ\xa04\x14\x1b\xb1\n\xa3\x13\xf4C\xe5;\xceg\xebV3\x9bV\x00{\xb0\x1e\xc3\x12\x037Rۏ]E\xf6E\x00\xfdq\x1aK\xdb\b\x9c7\xbd\x81y\x8d|5\x02\x98\x82\x87\x15\x05λ\xf9\x93e\xf2S\xfe\xcbz\x80Ax\x85\xbe\xff\x01 \\\xc6a]\x80\xa7]\xfe\xeb\x01.%\xa4\xfd\x11\xae\xe5\x007!\xc0S%n\xc9L\xfc\n\x83<\xc0Nā)\x92\xfb؊\au\n\xa3\xb0-[\xa9\x10{\x11\a\xbc\x060\v7ѣ\xc0\xc3e\xac\xc5\xd7&\x11\xd6\xe1d\"A\xfc\xfb\xc8\xc2\xe3\xae\f\x9cÊ&\xc6sp'\xc1\xf9QlN\xd0\v\x95C\xd8\x11\x00\xc3\xf1\x12}\n\f#p\x00\xd4\xcb\x1a\xac\xc6v<\xc9~\x9c\x98?`-@\xa2rN\b\x80-8\xd2D\xb9\xb6_u:\x11t\x12\xf2\xbf\xf7\xc6[\fN\xcc\xc2\xe6\x008\x8e8\xc9ER\x04\x10E*V\x1c\xa6w9\xe3\xe8\x03\x01\xd6/\x01\xe2T\x00\\\xc7\xfc&\xca'\xb2^\x90௦\x12\xdb\x1a\xdb\x11 \xf1\x8c5\xae\x81\x83\xbb\x01p\f\x17\xf1\x05\x9f\xb3\xe7\x00L\xc74\xcc\xc8>\x97\x01\xa8\xd7\xed\xd9\x00h6~\xb4jF\xed\x04mf\x1b۶;_\a\xaa\x0e4\x06c\x11ϑ\x18\x92[\xd1ޣ\xd9uJ\x95\x19X\x99\xcd\fQ듥*\x80Ve\xbc\b\xe8I\x15\x00\xbd\xb2\xbb\x1fi.+\x87\xab\x00\x88+\x1cW\xb9;\xb2\xa6\n\x80}\xcdF\xae&T17v\x96\xe2v\xe5^V/\xca\xfa\x89ι\xa1]\x80\xb8fo\xcaF\xce\xf4\xe7\xe2V\xbb\x001#\x9c\xe9\x06@\f\xbc1?\xb4]\a\xae`qI\x80\xe7\x18_E!Z\x84\xab%\x83\xffF\xb4운\xb3\x05e\xff\xfdkLŷ*\x00j\xcd$1\x031\t/\xab\x9b\x1b:M\xbb\x93\x81帐\x188\x86σ8\xd05\x86\xd7ە\x05H\r\x1e\xf3^\xbc\x9c\x9e\xcdޖ\vy\xcb\x00\x8c\xc8F\xe9\x18\xa7cM\xc8y\xfd\x89\x17\x88\x13\x1e\xafu\x11<\x06\x9c\x96\xf2\x0f\x00\xddl\xca\a\xee\xc0?\x00\x00\x00\x00IEND\xaeB`\x82"),
	}
	file7 := &embedded.EmbeddedFile{
		Filename:    "images/logo.svg",
		FileModTime: time.Unix(1591631914, 0),

		Content: string("<?xml version=\"1.0\" encoding=\"utf-8\"?>\r\n<!-- Generator: Adobe Illustrator 13.0.0, SVG Export Plug-In . SVG Version: 6.00 Build 14948)  -->\r\n<!DOCTYPE svg PUBLIC \"-//W3C//DTD SVG 1.1//EN\" \"http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd\">\r\n<svg version=\"1.1\" id=\"Layer_1\" xmlns=\"http://www.w3.org/2000/svg\" xmlns:xlink=\"http://www.w3.org/1999/xlink\" x=\"0px\" y=\"0px\"\r\n\t width=\"241.839px\" height=\"50.047px\" viewBox=\"0 0 241.839 50.047\" enable-background=\"new 0 0 241.839 50.047\"\r\n\t xml:space=\"preserve\">\r\n<g>\r\n\t<path fill=\"#FFFFFF\" d=\"M70.995,21.187c0.952,0,1.987,0.501,3.104,1.504l2.38-2.38c-1.713-1.714-3.587-2.571-5.618-2.571\r\n\t\tc-1.84,0-3.357,0.533-4.551,1.6c-1.193,1.066-1.79,2.412-1.79,4.037c0,1.079,0.261,1.993,0.781,2.742\r\n\t\tc0.521,0.711,1.454,1.44,2.799,2.189c0.762,0.406,1.365,0.755,1.81,1.047c0.443,0.292,0.755,0.521,0.933,0.687\r\n\t\tc0.368,0.33,0.551,0.762,0.551,1.295c0,0.608-0.218,1.104-0.656,1.484c-0.437,0.382-1.019,0.572-1.742,0.572\r\n\t\tc-1.283,0-2.488-0.68-3.619-2.038l-2.761,1.923c1.727,2.374,3.903,3.562,6.532,3.562c1.93,0,3.497-0.507,4.703-1.523\r\n\t\tc1.181-1.027,1.772-2.393,1.772-4.094c0-1.181-0.293-2.14-0.877-2.876c-0.545-0.749-1.625-1.574-3.236-2.475\r\n\t\tc-1.131-0.623-1.873-1.118-2.229-1.486c-0.355-0.355-0.533-0.793-0.533-1.314c0-0.533,0.213-0.981,0.639-1.342\r\n\t\tC69.81,21.368,70.348,21.187,70.995,21.187\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M95.037,20.691c-0.952-1.028-1.885-1.752-2.799-2.171c-0.902-0.419-1.994-0.628-3.275-0.628\r\n\t\tc-2.717,0-4.945,0.946-6.685,2.837c-1.713,1.892-2.571,4.317-2.571,7.275c0,2.563,0.75,4.651,2.248,6.265\r\n\t\tc1.498,1.612,3.434,2.418,5.808,2.418c1.93,0,3.796-0.787,5.598-2.361l-0.248,1.942h4.209l2.209-17.938h-4.209L95.037,20.691z\r\n\t\t M92.714,31.317c-1.13,1.257-2.47,1.885-4.019,1.885c-1.396,0-2.539-0.501-3.427-1.504c-0.889-1.003-1.332-2.292-1.332-3.865\r\n\t\tc0-1.804,0.533-3.327,1.599-4.57c1.067-1.232,2.387-1.847,3.961-1.847c1.447,0,2.627,0.501,3.541,1.504\r\n\t\tc0.914,0.99,1.371,2.279,1.371,3.866C94.408,28.524,93.844,30.035,92.714,31.317\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M113.992,17.892c-1.701,0-3.383,0.698-5.047,2.095l0.192-1.657h-4.19l-2.209,17.938h4.208l1.067-8.646\r\n\t\tc0.088-0.724,0.19-1.36,0.304-1.914c0.114-0.552,0.254-1.019,0.419-1.4c0.305-0.749,0.775-1.383,1.409-1.904\r\n\t\tc0.775-0.66,1.701-0.99,2.781-0.99c1.752,0,2.628,0.812,2.628,2.437c0,0.242-0.009,0.515-0.029,0.819\r\n\t\tc-0.019,0.305-0.054,0.654-0.105,1.048l-1.295,10.549h4.189l1.2-9.636c0.153-1.27,0.229-2.361,0.229-3.275\r\n\t\tc0-1.688-0.515-3.021-1.542-3.999C117.146,18.381,115.744,17.892,113.992,17.892\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M139.928,21.053c-1.053-1.155-2.063-1.974-3.027-2.457c-0.927-0.469-2.068-0.704-3.429-0.704\r\n\t\tc-2.651,0-4.854,0.933-6.606,2.799c-1.752,1.866-2.628,4.209-2.628,7.027c0,2.641,0.8,4.799,2.398,6.475\r\n\t\tc1.601,1.663,3.657,2.494,6.17,2.494c1.105,0,2.031-0.152,2.78-0.457c0.762-0.33,1.664-0.958,2.705-1.885l-0.229,1.923h4.189\r\n\t\tl3.961-32.163h-4.209L139.928,21.053z M137.604,31.241c-1.104,1.295-2.457,1.942-4.057,1.942c-1.535,0-2.767-0.477-3.693-1.428\r\n\t\tc-0.928-0.966-1.391-2.241-1.391-3.827c0-1.879,0.539-3.435,1.62-4.666c1.078-1.231,2.443-1.847,4.093-1.847\r\n\t\tc1.524,0,2.755,0.47,3.695,1.409c0.938,0.939,1.409,2.171,1.409,3.694C139.281,28.397,138.722,29.972,137.604,31.241\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M159.945,17.893c-2.007,0-3.917,0.774-5.732,2.323l0.229-1.885h-4.209l-3.467,28.202h4.208l1.581-12.835\r\n\t\tc1.029,1.131,1.999,1.911,2.913,2.344c0.877,0.431,2,0.646,3.371,0.646c2.679,0,4.92-0.958,6.723-2.875\r\n\t\tc1.814-1.93,2.724-4.304,2.724-7.123c0-2.539-0.793-4.64-2.381-6.303C164.319,18.724,162.332,17.893,159.945,17.893\r\n\t\t M162.439,31.355c-1.092,1.219-2.438,1.828-4.036,1.828c-1.524,0-2.769-0.495-3.733-1.485c-0.964-0.965-1.447-2.24-1.447-3.827\r\n\t\tc0-1.815,0.571-3.345,1.713-4.589c1.144-1.244,2.565-1.866,4.267-1.866c1.422,0,2.59,0.514,3.504,1.542\r\n\t\tc0.914,1.016,1.371,2.317,1.371,3.904C164.077,28.613,163.532,30.111,162.439,31.355\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M177.312,7.457c-0.608,0-1.136,0.222-1.581,0.667c-0.432,0.431-0.646,0.958-0.646,1.58\r\n\t\ts0.222,1.174,0.666,1.657c0.458,0.457,0.997,0.685,1.62,0.685c0.621,0,1.154-0.222,1.599-0.666\r\n\t\tc0.444-0.444,0.666-0.971,0.666-1.581c0-0.634-0.233-1.18-0.705-1.637C178.473,7.692,177.934,7.457,177.312,7.457\"/>\r\n\t<polygon fill=\"#FFFFFF\" points=\"171.903,36.269 176.111,36.269 178.32,18.331 174.112,18.331 \t\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M193.899,17.893c-2.007,0-3.917,0.774-5.733,2.323l0.229-1.885h-4.209l-3.466,28.202h4.208l1.581-12.835\r\n\t\tc1.028,1.131,1.999,1.911,2.913,2.344c0.876,0.431,1.999,0.646,3.37,0.646c2.68,0,4.92-0.958,6.723-2.875\r\n\t\tc1.815-1.93,2.724-4.304,2.724-7.123c0-2.539-0.793-4.64-2.38-6.303S196.286,17.893,193.899,17.893 M196.393,31.355\r\n\t\tc-1.092,1.219-2.438,1.828-4.036,1.828c-1.524,0-2.769-0.495-3.733-1.485c-0.964-0.965-1.447-2.24-1.447-3.827\r\n\t\tc0-1.815,0.571-3.345,1.714-4.589c1.144-1.244,2.565-1.866,4.267-1.866c1.422,0,2.59,0.514,3.504,1.542\r\n\t\tc0.914,1.016,1.371,2.317,1.371,3.904C198.031,28.613,197.485,30.111,196.393,31.355\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M237.901,17.854c-1.446,0-2.883,0.673-4.303,2.019l0.19-1.543h-4.209l-2.209,17.939h4.209l1.104-8.951\r\n\t\tc0.113-0.99,0.271-1.84,0.465-2.551c0.197-0.711,0.445-1.292,0.744-1.743c0.297-0.45,0.651-0.781,1.058-0.99\r\n\t\tc0.405-0.21,0.875-0.314,1.407-0.314c0.775,0,1.55,0.33,2.324,0.99l2.687-3.466C240.326,18.317,239.172,17.854,237.901,17.854\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M25.326,0.956c-12.765,0-23.262,9.864-24.287,22.369c0.73-0.245,1.472-0.456,2.229-0.605\r\n\t\tc3.481-0.681,7.19-0.283,10.54,0.837c0.207,0.069,0.408,0.161,0.613,0.237c-0.498-1.958-0.792-4.038-0.408-5.984\r\n\t\tc0.965-4.876,6.625-8.221,11.35-6.719c2.53,0.804,4.728,2.675,5.687,5.161c0.523,1.356,0.817,2.831,0.963,4.289\r\n\t\tc4.48,1.028,8.959,2.057,13.44,3.085c1.019,0.234,0.695,1.644-0.232,1.716c-2.56,0.198-5.118,0.441-7.667,0.73\r\n\t\tc-1.786,0.203-3.575,0.421-5.339,0.763c-1.936,0.375-3.769,1.023-4.754,2.867c-0.789,1.477-1.026,3.645-1.493,5.34\r\n\t\tc-0.76,2.755-1.751,5.482-3.222,7.943c-1.253,2.097-2.912,4.095-4.931,5.536c2.367,0.769,4.891,1.188,7.511,1.188\r\n\t\tc13.441,0,24.377-10.936,24.377-24.376C49.703,11.892,38.767,0.956,25.326,0.956\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M24.428,20.29c1.057,0,1.913-0.857,1.913-1.914c0-1.056-0.856-1.913-1.913-1.913s-1.913,0.857-1.913,1.913\r\n\t\tC22.515,19.433,23.371,20.29,24.428,20.29\"/>\r\n\t<path fill=\"#FFFFFF\" d=\"M210.181,28.805c0,1.319,0.457,2.404,1.371,3.256c0.927,0.851,2.106,1.275,3.542,1.275\r\n\t\tc1.968,0,3.478-0.735,4.531-2.209l3.144,1.733c-1.092,1.472-2.216,2.501-3.371,3.085c-1.168,0.596-2.628,0.895-4.38,0.895\r\n\t\tc-2.781,0-4.996-0.818-6.646-2.457c-1.649-1.637-2.476-3.826-2.476-6.569c0-2.869,0.908-5.268,2.723-7.198\r\n\t\tc1.804-1.917,4.081-2.875,6.837-2.875c2.666,0,4.768,0.863,6.304,2.59c1.561,1.739,2.342,4.068,2.342,6.988\r\n\t\tc0,0.305-0.02,0.8-0.058,1.486H210.181z M219.834,25.414c-0.442-2.768-1.974-4.151-4.589-4.151c-2.475,0-4.088,1.383-4.837,4.151\r\n\t\tH219.834z\"/>\r\n</g>\r\n</svg>\r\n"),
	}
	file8 := &embedded.EmbeddedFile{
		Filename:    "index.html",
		FileModTime: time.Unix(1591639927, 0),

		Content: string("<!DOCTYPE html>\r\n<html lang=\"en\">\r\n  <head>\r\n    <meta charset=\"utf-8\" />\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1, shrink-to-fit=no\" />\r\n    <title>Request Sandpiper Access</title>\r\n    <link rel=\"icon\" type=\"image/png\" href=\"images/favicon-32x32.png\" sizes=\"32x32\" />\r\n    <link rel=\"icon\" type=\"image/png\" href=\"images/favicon-16x16.png\" sizes=\"16x16\" />\r\n    <link rel=\"stylesheet\" href=\"css/style.css\" />\r\n  </head>\r\n  <div class=\"register-page\">\r\n    <a href=\"https://sandpiperframework.org\"><img src=\"images/logo.svg\" alt=\"logo\" class=\"logo\"/></a>\r\n    <div class=\"form\">\r\n      <form class=\"login-form\">\r\n        <input type=\"text\" placeholder=\"name\" />\r\n        <input type=\"text\" placeholder=\"email\" />\r\n        <input type=\"text\" placeholder=\"company\" />\r\n        <input type=\"text\" placeholder=\"sandpiper server id\" />\r\n        <select id = \"kind\">\r\n          <option value = \"1\">classification</option>\r\n          <option value = \"1\">Distributor</option>\r\n          <option value = \"2\">Retailer</option>\r\n          <option value = \"3\">Electronic Catalog</option>\r\n          <option value = \"4\">Other</option>\r\n        </select>\r\n        <button>register</button>\r\n        <p class=\"message\"><a href=\"#\">Terms & Conditions</a></p>\r\n      </form>\r\n    </div>\r\n  </div>\r\n</html>\r\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1591639927, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file8, // "index.html"

		},
	}
	dir2 := &embedded.EmbeddedDir{
		Filename:   "css",
		DirModTime: time.Unix(1591640285, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file3, // "css/style.css"

		},
	}
	dir4 := &embedded.EmbeddedDir{
		Filename:   "images",
		DirModTime: time.Unix(1591633059, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file5, // "images/favicon-16x16.png"
			file6, // "images/favicon-32x32.png"
			file7, // "images/logo.svg"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{
		dir2, // "css"
		dir4, // "images"

	}
	dir2.ChildDirs = []*embedded.EmbeddedDir{}
	dir4.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`../../../public`, &embedded.EmbeddedBox{
		Name: `../../../public`,
		Time: time.Unix(1591639927, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"":       dir1,
			"css":    dir2,
			"images": dir4,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"css/style.css":            file3,
			"images/favicon-16x16.png": file5,
			"images/favicon-32x32.png": file6,
			"images/logo.svg":          file7,
			"index.html":               file8,
		},
	})
}

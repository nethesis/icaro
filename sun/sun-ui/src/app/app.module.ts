import { HomeComponent } from './home/home.component';
import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http';
import { HttpModule  } from '@angular/http/';
import {routing} from './app.routing';
import { AppComponent }  from './app.component';
import {AuthenticationService} from './_services/authentication.service';
import {LoginComponent} from './login/login.component';
import {NavbarComponent} from './navbar/navbar.component';
import {SidenavComponent} from './sidenav/sidenav.component';
import {NotFoundComponent} from './404/404.component';


@NgModule({
	imports:
	[
		BrowserModule,
		FormsModule,
		HttpClientModule,
		routing,
		HttpModule
	],

	declarations:
	[
		AppComponent,
		LoginComponent,
		HomeComponent,
		NavbarComponent,
		SidenavComponent,
		NotFoundComponent,
	],

	providers:
	[
		AuthenticationService
	],

	bootstrap:
	[
		AppComponent
	]
})
export class AppModule { }

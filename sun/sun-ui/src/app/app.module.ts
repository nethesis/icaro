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

@NgModule({
  imports:      [ BrowserModule, FormsModule, HttpClientModule, routing,  HttpModule ],
  declarations: [ AppComponent, LoginComponent, HomeComponent],
  providers:    [ AuthenticationService],
  bootstrap:    [ AppComponent ]
})
export class AppModule { }

import { NgModule }      from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { FormsModule } from '@angular/forms';
import { HttpClientModule, HTTP_INTERCEPTORS } from '@angular/common/http';

import {routing} from './app.routing'
import { AppComponent }  from './app.component';
import {AuthenticationService} from './_services/authentication.service';
import {LoginComponent} from './login/login.component';
import { HttpClient } from '@angular/common/http/src/client';

@NgModule({
  imports:      [ BrowserModule,FormsModule, HttpClientModule,routing],
  declarations: [ AppComponent, LoginComponent ],
  providers:    [AuthenticationService],
  bootstrap:    [ AppComponent ]
})
export class AppModule { }

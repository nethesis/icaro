import { HomeComponent } from './../home/home.component';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginComponent } from '../login/login.component';

/**
 * routes Constant for routing throw the pages, this is imported in app.modules.ts!
 */
const routes: Routes = [
  { path: '', redirectTo: '/login', pathMatch: 'full'},
  { path: 'login', component: LoginComponent},
  { path: 'home', component: HomeComponent}
];

export const routing = RouterModule.forRoot(routes);
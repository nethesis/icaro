import { ProfileComponent } from './../profile/profile.component';
import { AuthenticationGuard } from './../guards/authentication.guard';
import { HomeComponent } from './../home/home.component';
import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { LoginComponent } from '../login/login.component';

/**
 * routes Constant for routing throw the pages, this is imported in app.modules.ts!
 */
const routes: Routes = [
  {
    path: 'home',
    component: HomeComponent,
    canActivate: [AuthenticationGuard],
    children: [
      {
        path: 'profile',
        component: ProfileComponent,
        canActivate: [AuthenticationGuard]
      }
    ]
  },

  { path: 'login', component: LoginComponent },
  // otherwise redirect to home
  { path: '**', redirectTo: '/home', pathMatch: 'full' }
];

export const routing = RouterModule.forRoot(routes);

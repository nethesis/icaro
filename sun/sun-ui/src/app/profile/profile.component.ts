import { NgForm } from '@angular/forms';
import { ProfileService } from './../services/profile.service';
import { Component, OnInit, ViewEncapsulation } from '@angular/core';
import { CardConfig, InfoStatusCardConfig } from 'patternfly-ng';
import { Account, AccountFull } from '../interfaces/account';
import { error } from 'util';
import { Router, ActivatedRoute } from '@angular/router';

@Component({
  encapsulation: ViewEncapsulation.None,
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {
  config: CardConfig;
  model: any = {};
  accountArray: any[] = [];
  account: any = {};
  fullAccount: any = {};
  preferenceArray: any[] = [];
  passwordLength: number;
  apiResponse: string;

  constructor(
    private profileService: ProfileService,
    private activateroute: ActivatedRoute,
    private router: Router
  ) {
    this.passwordLength = 0;
  }

  ngOnInit() {
    this.config = {
      noPadding: false,
      topBorder: true
    } as CardConfig;
    this.getAccount();
    this.getPreferences();
  }

  private postAccount() {
    this.fullAccount.email = 'test@hotmail.com';
    this.fullAccount.name = 'Test';
    this.fullAccount.username = 'Test';
    this.fullAccount.uuid = '29';
    this.fullAccount.hotspotid = '29';
    this.fullAccount.password = 'password';
    this.fullAccount.type = 'type1';

    this.profileService.postAccount(this.fullAccount).subscribe(
      data => {
        console.log(data);
      },
      error => {
        console.log(error);
      }
    );
  }

  private getAccount() {
    this.profileService.getAccount(5).subscribe(
      data => {
        this.accountArray = data;
        console.log(data);
      },
      error => {
        console.log(error.json());
      }
    );
  }

  private getPreferences() {
    this.profileService.getPreferences().subscribe(
      data => {
        let myArray = [];
        for (let key in error) {
          myArray.push(data[key]);
        }
        this.preferenceArray = myArray;
      },
      error => {}
    );
  }
}

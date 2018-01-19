import { NgForm } from '@angular/forms';
import { ProfileService } from '../../services/profile.service';
import { Component, OnInit, ViewEncapsulation, ViewChild, ElementRef} from '@angular/core';
import { error } from 'util';
import * as jQuery from 'jquery';

declare var $: any;
@Component({
  encapsulation: ViewEncapsulation.None,
  selector: 'app-change-password-modal',
  templateUrl: './change-password-modal.component.html',
  styleUrls: ['./change-password-modal.component.css']
})
export class ChangePasswordModalComponent implements OnInit {
  @ViewChild('myModal') myModal: ElementRef;
  private id_profile = Number(localStorage.getItem('id'));
  model: any = {};
  items: any[] = [];
  account: any = {};
  fullAccount: any = {};
  preferenceArray: any[] = [];
  passwordLength = 0;
  apiResponse: string;

  constructor(private profileService: ProfileService) {}

  ngOnInit() {
     this.getAccount(this.id_profile);
  }

  private getPasswordLength(value: string = '') {
    this.passwordLength = value.length;
  }

  private changePassword(changepasswordform: NgForm) {
   this.account.email = this.fullAccount.email;
    this.account.name = this.fullAccount.name;
    this.account.username = this.fullAccount.username;
    this.account.password = this.model.newPassword;
    this.account.id = this.fullAccount.id;
    this.profileService.changePassword(this.account).subscribe(
      data => {
        this.apiResponse = data.status;
        changepasswordform.reset();
        this.passwordLength = 0;
        if (this.apiResponse === 'success') {
          $(this.myModal.nativeElement).modal('hide');
        }
      },
      error => {
        this.apiResponse = 'error';
        console.log(error);
      }
    );
  }

  private getAccount(id:number) {
    this.profileService.getProfilebyId(id).subscribe(
      data => {
        this.fullAccount = data;
      },
      error => {
        console.log(error);
      }
    );
  }


  private resetFields(changepasswordform: NgForm) {
    changepasswordform.reset();
    this.passwordLength = 0;
    this.apiResponse = '';
  }
}

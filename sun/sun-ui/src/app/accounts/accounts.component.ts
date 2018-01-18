import { error } from 'util';
import { ProfileService } from './../services/profile.service';
import { Component, OnInit, ViewChild, ElementRef } from '@angular/core';
import {
  ActionConfig,
  ListConfig,
  Action,
  ListEvent,
  PaginationConfig,
  ToolbarConfig
} from 'patternfly-ng';
import { PaginationEvent } from 'patternfly-ng/dist/src/app/pagination/pagination.module';
import { OnChanges } from '@angular/core/src/metadata/lifecycle_hooks';
import * as jQuery from 'jquery';

declare var $: any;
@Component({
  selector: 'app-accounts',
  templateUrl: './accounts.component.html',
  styleUrls: ['./accounts.component.css']
})
export class AccountsComponent implements OnInit {
  @ViewChild('myModal') myModal: ElementRef;

  itemToEdit: any[] = [];
  seletedItem: any[] = [];
  isDataArrived = false;
  accountList: any[];
  actionConfig: ActionConfig;
  actionConfig1: ActionConfig;
  actionsText: string = '';
  items: any[];
  listConfig: ListConfig;
  paginationConfig: PaginationConfig;

  constructor(private profileService: ProfileService) {}
  ngOnInit() {
    this.getAllAccounts();
    this.actionConfig = {
      primaryActions: [
        {
          id: 'Edit',
          title: 'Edit',
          tooltip: 'Edit Accounts'
        },
        {
          id: 'action2',
          title: 'Delete',
          tooltip: 'Delete account'
        }
      ]
    } as ActionConfig;

    this.listConfig = {
      dblClick: false,
      multiSelect: false,
      selectItems: false,
      selectionMatchProp: 'name',
      showCheckbox: true,
      useExpandItems: false
    } as ListConfig;

    this.paginationConfig = {
      pageNumber: 1,
      pageSize: 3,
      pageSizeIncrements: [2, 3, 4],
      totalItems: this.accountList.length
    } as PaginationConfig;

    this.updateItems();
  }

  private getAllAccounts() {
    this.profileService.getAllAccounts().subscribe(
      data => {
        this.isDataArrived = true;
        this.accountList = data;
        console.log(this.accountList);
      },
      error => {
        console.log(error);
      }
    );
  }

  // Actions

  handleAction($event: Action, item: any): void {
    if ($event.id === 'Edit') {
      this.itemToEdit = item;
      $(this.myModal.nativeElement).modal('show');
    }
  }

  handleSelectionChange($event: ListEvent): void {
    if($event.selectedItems.length > 0){
      this.seletedItem = $event.selectedItems;
    }else{
      this.seletedItem = [];
    }
  }

  handleClick($event: ListEvent, item: any): void {
  }

  handlePageSize($event: PaginationEvent) {
    this.updateItems();
  }

  handlePageNumber($event: PaginationEvent) {
    this.updateItems();
  }

  // Pagination
  updateItems() {
    this.items = this.accountList
      .slice(
        (this.paginationConfig.pageNumber - 1) * this.paginationConfig.pageSize,
        this.paginationConfig.totalItems
      )
      .slice(0, this.paginationConfig.pageSize);
  }
}

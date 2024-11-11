import { Injectable } from '@angular/core';
import { MatTableDataSource } from '@angular/material/table';
import { UserClass } from '../../../api/models/userClass';
import { AdminApiService } from '../../../api/admin/admin-api.service';
import { MatPaginator } from '@angular/material/paginator';
import { NonNullableFormBuilder } from '@angular/forms';

@Injectable({
  providedIn: 'root'
})
export class UsersTableService {
  public dataSource = new MatTableDataSource<UserClass>([]);
  public conditions = this.fb.group({
    roleIds: [[] as Array<number>],
  })
  private paginator!: MatPaginator

  constructor(
    private adminApi: AdminApiService,
    private fb: NonNullableFormBuilder,
  ) {

  }

  public setPaginator(p: MatPaginator) {
    this.paginator = p
    this.paginator.page.subscribe(() => {
      this.refreshTable();
    });
    this.conditions.valueChanges.subscribe(() => {
      this.paginator.pageIndex = 0
      this.refreshTable()
    })
  }

  refreshTable() {
    console.log('refresh table', this.conditions.value)
    this.adminApi
      .getUsers({
        desc: true,
        keyword: '',
        limit: this.paginator.pageSize,
        page: this.paginator.pageIndex,
        roleIds: this.conditions.value.roleIds || [],
      })
      .subscribe((data) => {
        this.dataSource.data = data.users;
        this.paginator.pageSize = data.limit;
        this.paginator.pageIndex = data.page;
        this.paginator.length = data.count;
      });
  }
}

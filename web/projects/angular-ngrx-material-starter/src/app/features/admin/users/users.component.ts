import { AfterViewInit, ChangeDetectionStrategy, Component, OnInit, ViewChild } from '@angular/core';
import { MatTableDataSource } from '@angular/material/table';
import { User } from '../../../api/models/user';
import { ApiService } from '../../../api/api.service';
import { MatPaginator } from '@angular/material/paginator';

@Component({
  selector: 'anms-users',
  templateUrl: './users.component.html',
  styleUrls: ['./users.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class UsersComponent implements OnInit, AfterViewInit {
  @ViewChild(MatPaginator) paginator: MatPaginator | undefined;

  public displayedColumns: string[] = ['actions', 'id', 'avatar', 'nickname', 'phone', 'roles'];
  public dataSource = new MatTableDataSource<User>([]);

  constructor(
    private api: ApiService,
  ) {
  }

  ngOnInit(): void {
  }

  refreshTable() {
    if (this.paginator === undefined) {
      return
    }
    this.api.getUsers({
      desc: true,
      keyword: '',
      limit: this.paginator.pageSize,
      page: this.paginator.pageIndex
    }).subscribe(
      data => {
        this.dataSource.data = data.users
        if (this.paginator === undefined) {
          return
        }
        this.paginator.pageSize = data.limit
        this.paginator.pageIndex = data.page
        this.paginator.length = data.count
      }
    )
  }

  ngAfterViewInit(): void {
    if (this.paginator === undefined) {
      return
    }
    this.paginator.pageSize = 10
    this.paginator.page.subscribe(
      (v) => {
        this.refreshTable()
      }
    )
    // this.dataSource.paginator = this.paginator
    this.refreshTable()
  }
}

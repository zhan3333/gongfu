import { ChangeDetectionStrategy, Component, Inject, OnInit } from '@angular/core';
import {
  FormArray,
  FormControl,
  FormGroup,
  FormsModule,
  NonNullableFormBuilder,
  ReactiveFormsModule,
  Validators
} from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { NotificationService } from '../../../core/notifications/notification.service';
import { AdminApiService, MemberCourse } from '../../../api/admin/admin-api.service';
import { MatChipInputEvent, MatChipsModule } from '@angular/material/chips';
import { displayRoleName, ROLE_COACH, ROLE_MEMBER, User } from '../../../api/models/user';
import { Levels } from '../../../services/coach-level';
import { faBan, faChalkboardTeacher, faEdit, faFolder, faGraduationCap } from '@fortawesome/free-solid-svg-icons';
import { MatButtonModule } from '@angular/material/button';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';
import { MatIconModule } from '@angular/material/icon';
import { MatOptionModule, MatRippleModule } from '@angular/material/core';
import { DatePipe, NgFor, NgIf } from '@angular/common';
import { MatSelectModule } from '@angular/material/select';
import { MatInputModule } from '@angular/material/input';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatCardModule } from '@angular/material/card';
import { MAT_DIALOG_DATA, MatDialog, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { MatListModule } from '@angular/material/list';
import { TeachingRecordDialog } from './teaching-dialog';
import { StudyRecordDialog } from './study-dialog';
import { StudyRecord } from '../../../api/models/study-record';
import { TeachingRecord } from '../../../api/models/teaching-record';
import { MatBottomSheet } from '@angular/material/bottom-sheet';
import { BottomSheetComponent } from '../../../shared/bottom-sheet.component';
import { MatDatepickerModule } from '@angular/material/datepicker';

// 用户编辑页，可以设置用户角色，设置教练信息等
@Component({
  selector: 'anms-user-edit',
  templateUrl: './user-edit.component.html',
  changeDetection: ChangeDetectionStrategy.Default,
  standalone: true,
  imports: [
    MatCardModule,
    FormsModule,
    ReactiveFormsModule,
    MatFormFieldModule,
    MatInputModule,
    MatSelectModule,
    NgFor,
    MatOptionModule,
    MatChipsModule,
    MatIconModule,
    FontAwesomeModule,
    MatButtonModule,
    MatListModule,
    NgIf,
    MatDialogModule,
    MatRippleModule,
    DatePipe
  ]
})
export class UserEditComponent implements OnInit {
  public faBan = faBan;
  public id = 0;

  public form = new FormGroup({
    phone: new FormControl('', {nonNullable: true}),
    nickname: new FormControl('', {nonNullable: true}),
    level: new FormControl('', {nonNullable: true}),
    teachingSpace: new FormControl('', {nonNullable: true}),
    teachingAge: new FormControl('', {nonNullable: true}),
    teachingExperiences: new FormControl<string[]>([], {nonNullable: true, validators: [Validators.required]}),
    roleNames: new FormControl<string[]>([], {nonNullable: true, validators: [Validators.required]}),
    // 会员类型
    memberType: new FormControl('', {nonNullable: true}),
    memberCourses: new FormArray([new FormGroup({
      // 课程名称
      name: new FormControl('', {nonNullable: true}),
      // 课程开始日期
      startDate: new FormControl('', {nonNullable: true}),
      // 课程结束日期
      endDate: new FormControl('', {nonNullable: true}),
      // 总课时数
      total: new FormControl(0, {nonNullable: true}),
      // 剩余课时数
      remain: new FormControl(0, {nonNullable: true}),
    })])
  });
  public loading = false;
  // 角色名选项
  public roleNames: string[] = [];
  public displayRoleName = displayRoleName;

  public Levels = Levels;
  public user: User | null = null;
  public memberCourses: MemberCourse[] = [];
  protected readonly faEdit = faEdit;
  protected readonly faFolder = faFolder;
  protected readonly faChalkboardTeacher = faChalkboardTeacher;
  protected readonly faGraduationCap = faGraduationCap;

  constructor(
    public adminApi: AdminApiService,
    public route: ActivatedRoute,
    public notification: NotificationService,
    public router: Router,
    private dialog: MatDialog,
    private _bottomSheet: MatBottomSheet,
  ) {
  }

  // 是否勾选了教练选项
  get selectedCoach() {
    return this.form.controls.roleNames.value.indexOf(ROLE_COACH) !== -1
  }

  // 是否勾选了会员角色
  get selectedMember() {
    return this.form.controls.roleNames.value.indexOf(ROLE_MEMBER) !== -1
  }

  ngOnInit(): void {
    this.id = parseInt(this.route.snapshot.params['id'], 10);
    if (this.id === 0) {
      this.notification.error('invalid user id');
      return;
    }
    this.refresh();
    this.refreshMemberCourses()
  }

  refresh() {
    this.adminApi.getUser(this.id).subscribe((user) => {
      this.user = user;
      this.form.patchValue({
        nickname: user.nickname,
        phone: user.phone,
        roleNames: user.roleNames,
      });
    });
    this.adminApi.getCoach(this.id).subscribe((coach) => {
      this.form.patchValue({
        level: coach.level,
        teachingSpace: coach.teachingSpace,
        teachingAge: coach.teachingAge,
        teachingExperiences: coach.teachingExperiences || [],
      });
    });
    this.adminApi.getRoleNames().subscribe((roleNames) => {
      this.roleNames = roleNames;
    });
  }

  onClickSave() {
    this.loading = true;
    this.adminApi
      .updateUser(this.id, {
        level: this.form.controls.level.value,
        nickname: this.form.controls.nickname.value,
        phone: this.form.controls.phone.value,
        teachingAge: this.form.controls.teachingAge.value,
        teachingExperiences: this.form.controls.teachingExperiences.value,
        teachingSpace: this.form.controls.teachingSpace.value,
        roleNames: this.form.controls.roleNames.value
      })
      .subscribe({
        next: () => {
          this.notification.success('保存成功');
        },
        error: (error) => {
          this.notification.error('保存失败: ' + error);
        }
      })
      .add(() => {
        this.loading = false;
      });
  }

  onRemoveExp(exp: any) {
    const exps: string[] = this.form.controls.teachingExperiences.value;
    const newExps: string[] = [];
    for (const item of exps) {
      if (exp === item) {
        continue;
      }
      newExps.push(item);
    }
    this.form.patchValue({
      teachingExperiences: [...newExps]
    });
  }

  onAddExp($event: MatChipInputEvent) {
    const value = ($event.value || '').trim();

    // Add our fruit
    if (value) {
      const exps: string[] = this.form.controls.teachingExperiences.value
      exps.push(value);
      this.form.patchValue({
        teachingExperiences: [...exps]
      });
    }

    // Clear the input value
    $event.chipInput?.clear();
  }

  // 打开新增授课记录弹窗
  openTeachingRecordDialog(record?: TeachingRecord) {
    this.dialog
      .open(TeachingRecordDialog, {data: record})
      .afterClosed()
      .subscribe(
        (result?: {
          data?: { id: number; date: string; address: string };
          isDelete?: boolean;
        }) => {
          if (!result) {
            return;
          }
          if (result.isDelete) {
            // 删除
          }
          if (result.data) {
            this.adminApi
              .editTeachingRecord({
                id: result.data.id,
                date: result.data.date,
                address: result.data.address,
                userId: this.user!.id
              })
              .subscribe({
                next: () => {
                  this.refresh();
                  this.notification.success('ok');
                }
              });
          }
        }
      );
  }

  // 打开新增授课记录弹窗
  openStudyRecordDialog(record?: StudyRecord) {
    this.dialog
      .open(StudyRecordDialog, {data: record})
      .afterClosed()
      .subscribe(
        (result?: {
          data?: { id: number; date: string; content: string };
          isDelete?: boolean;
        }) => {
          if (!result) {
            return;
          }
          if (result.isDelete) {
            this.adminApi.deleteStudyRecord(record!.id).subscribe({
              next: () => {
                this.refresh();
                this.notification.success('ok');
              }
            });
          }
          if (result.data) {
            this.adminApi
              .editStudyRecord({
                id: result.data.id,
                date: result.data.date,
                content: result.data.content,
                userId: this.user!.id
              })
              .subscribe({
                next: () => {
                  this.refresh();
                  this.notification.success('ok');
                }
              });
          }
        }
      );
  }

  // open confirm dialog, delete memberCourse
  deleteCourse(course: MemberCourse) {
    this._bottomSheet.open(BottomSheetComponent, {
      data: new Map<string, string>([
        ['ok', '删除'],
        ['cancel', '取消']
      ])
    }).afterDismissed().subscribe(v => {
      if (v === 'ok') {
        this.adminApi.delMemberCourse(course.id).subscribe({
          next: () => {
            this.refreshMemberCourses();
            this.notification.success('ok');
          }
        });
      }
    })
  }

  // open dialog create member course
  openCreateMemberCourseDialog() {
    this.dialog
      .open(CreateMemberCourseDialog, {})
      .afterClosed()
      .subscribe((result) => {
        if (result) {
          this.adminApi
            .createMemberCourse({
              userId: this.id,
              ...result
            })
            .subscribe({
              next: () => {
                this.refreshMemberCourses();
                this.notification.success('ok');
              }
            });
        }
      });
  }

  openEditMemberCourseDialog(course: MemberCourse) {
    this.dialog
      .open(UpdateMemberCourseDialog, {data: course})
      .afterClosed()
      .subscribe((result) => {
        if (result) {
          this.adminApi
            .updateMemberCourse(course.id, {
              ...result
            })
            .subscribe({
              next: () => {
                this.refreshMemberCourses();
                this.notification.success('ok');
              }
            });
        }
      });
  }

  private refreshMemberCourses() {
    this.adminApi.getMemberCourses(this.id).subscribe((courses) => {
      this.memberCourses = courses;
    });
  }
}

@Component({
  selector: 'create-member-course-dialog',
  template: `
    <h2 mat-dialog-title>创建会员课程</h2>
    <mat-dialog-content>
      <form [formGroup]="form">
        <mat-form-field>
          <mat-label>课程名称</mat-label>
          <input matInput formControlName="name"/>
        </mat-form-field>
        <mat-form-field>
          <mat-label>开始日期</mat-label>
          <input matInput [matDatepicker]="startDatePicker" formControlName="startTime"/>
          <mat-datepicker-toggle matSuffix [for]="startDatePicker"></mat-datepicker-toggle>
          <mat-datepicker #startDatePicker></mat-datepicker>
        </mat-form-field>
        <mat-form-field>
          <mat-label>结束日期</mat-label>
          <input matInput [matDatepicker]="endDatePicker" formControlName="endTime"/>
          <mat-datepicker-toggle matSuffix [for]="endDatePicker"></mat-datepicker-toggle>
          <mat-datepicker #endDatePicker></mat-datepicker>
        </mat-form-field>
        <mat-form-field>
          <mat-label>总课时数</mat-label>
          <input matInput type="number" formControlName="total"/>
        </mat-form-field>
        <mat-form-field>
          <mat-label>备注</mat-label>
          <input matInput formControlName="remark"/>
        </mat-form-field>
      </form>
    </mat-dialog-content>
    <mat-dialog-actions>
      <button mat-button (click)="onNoClick()">取消</button>
      <button mat-raised-button color="primary" [mat-dialog-close]="form.value" cdkFocusInitial>确定</button>
    </mat-dialog-actions>
  `,
  standalone: true,
  imports: [
    MatFormFieldModule,
    MatInputModule,
    FormsModule,
    MatButtonModule,
    ReactiveFormsModule,
    MatDatepickerModule,
    MatDialogModule,
  ],
})
export class CreateMemberCourseDialog {
  public form = this.fb.group({
    name: ['', Validators.required],
    startTime: ['', Validators.required],
    endTime: ['', Validators.required],
    total: [0, Validators.required],
    remark: [''],
  })

  constructor(
    public dialogRef: MatDialogRef<CreateMemberCourseDialog>,
    private fb: NonNullableFormBuilder,
  ) {
    this.form.valueChanges.subscribe(v => console.log('vvv', v))
  }

  onNoClick(): void {
    this.dialogRef.close();
  }
}


@Component({
  selector: 'update-member-course-dialog',
  template: `
    <h2 mat-dialog-title>编辑会员课程</h2>
    <mat-dialog-content>
      <form [formGroup]="form">
        <mat-form-field>
          <mat-label>课程名称</mat-label>
          <input matInput formControlName="name"/>
        </mat-form-field>
        <mat-form-field>
          <mat-label>开始日期</mat-label>
          <input matInput [matDatepicker]="startDatePicker" formControlName="startTime"/>
          <mat-datepicker-toggle matSuffix [for]="startDatePicker"></mat-datepicker-toggle>
          <mat-datepicker #startDatePicker></mat-datepicker>
        </mat-form-field>
        <mat-form-field>
          <mat-label>结束日期</mat-label>
          <input matInput [matDatepicker]="endDatePicker" formControlName="endTime"/>
          <mat-datepicker-toggle matSuffix [for]="endDatePicker"></mat-datepicker-toggle>
          <mat-datepicker #endDatePicker></mat-datepicker>
        </mat-form-field>
        <mat-form-field>
          <mat-label>总课时数</mat-label>
          <input matInput type="number" formControlName="total"/>
        </mat-form-field>
        <mat-form-field>
          <mat-label>剩余课时数</mat-label>
          <input matInput type="number" formControlName="remain"/>
        </mat-form-field>
        <mat-form-field>
          <mat-label>备注</mat-label>
          <input matInput formControlName="remark"/>
        </mat-form-field>
      </form>
    </mat-dialog-content>
    <mat-dialog-actions>
      <button mat-button (click)="onNoClick()">取消</button>
      <button mat-raised-button color="primary" [mat-dialog-close]="form.value" cdkFocusInitial>确定</button>
    </mat-dialog-actions>
  `,
  standalone: true,
  imports: [
    MatFormFieldModule,
    MatInputModule,
    FormsModule,
    MatButtonModule,
    ReactiveFormsModule,
    MatDatepickerModule,
    MatDialogModule,
  ],
})
export class UpdateMemberCourseDialog {
  public form = this.fb.group({
    name: ['', Validators.required],
    startTime: ['', Validators.required],
    endTime: ['', Validators.required],
    total: [0, Validators.required],
    remain: [0, Validators.required],
    remark: [''],
    status: [''],
  })

  constructor(
    public dialogRef: MatDialogRef<UpdateMemberCourseDialog>,
    @Inject(MAT_DIALOG_DATA) public data: MemberCourse,
    private fb: NonNullableFormBuilder,
  ) {
    this.form.patchValue(data)
  }

  onNoClick(): void {
    this.dialogRef.close();
  }
}

%global.counter = type { %global.counter.vtable.type*, i32, i32, i8*, void (i8*)* }
%global.counter.vtable.type = type { i8* ()*, void (i8*)*, void (i8*)*, void (i8*)*, void (i8*)*, void (i8*)* }
%global.base_class = type { %global.base_class.vtable.type* }
%global.base_class.vtable.type = type { i8* ()*, void (i8*)*, void (i8*)* }
%global.derive_class = type { %global.derive_class.vtable.type* }
%global.derive_class.vtable.type = type { i8* ()*, void (i8*)*, void (i8*)* }

@global.counter.vtable.data = global %global.counter.vtable.type { i8* ()* @global.counter.create, void (i8*)* @global.counter.destroy, void (i8*)* @global.counter.retain_shared, void (i8*)* @global.counter.release_shared, void (i8*)* @global.counter.retain_weak, void (i8*)* @global.counter.release_weak }
@global.base_class.vtable.data = global %global.base_class.vtable.type { i8* ()* @global.base_class.create, void (i8*)* @global.base_class.destroy, void (i8*)* @global.base_class.echo }
@global.derive_class.vtable.data = global %global.derive_class.vtable.type { i8* ()* @global.derive_class.create, void (i8*)* @global.derive_class.destroy, void (i8*)* @global.derive_class.echo }
@string.5bdaebb122965539cdd6ce77f212b65e = constant [15 x i8] c"create counter\00"
@string.f8f86b3941cca26e8c147322b9a8309f = constant [16 x i8] c"destroy counter\00"
@string.c1432ab71496ebb8b3e30bbcf37605e7 = constant [20 x i8] c"retain shared: %d \0A\00"
@string.ce99a7174da84f1767b9d5235ea3c24f = constant [21 x i8] c"release shared: %d \0A\00"
@string.21c67ac9191c65481dbab306227b4840 = constant [17 x i8] c"free object %p \0A\00"
@string.4fc1bf1a9ddd2be568f08ffc8ed6b9f0 = constant [18 x i8] c"free counter %p \0A\00"
@string.8d9c52192bdfa908703a004b070ff63e = constant [18 x i8] c"retain weak: %d \0A\00"
@string.ad2ae91c8542c824c6842efed6523f49 = constant [19 x i8] c"release weak: %d \0A\00"
@string.726bd3560bd4c136648f7760895d8d62 = constant [18 x i8] c"base construction\00"
@string.362aeeddb3d01da539cb6755bde46953 = constant [17 x i8] c"base destruction\00"
@string.9bcbb503bda6c8ad83f772846b706f08 = constant [13 x i8] c"echo in base\00"
@string.33b7808bf372c3d58730520160cb2c15 = constant [20 x i8] c"derive construction\00"
@string.ef25b0542457581e67c27a0dddb7bda5 = constant [19 x i8] c"derive destruction\00"
@string.895758554639f423e017c6610cbf460b = constant [15 x i8] c"echo in derive\00"

declare i32 @puts(i8* %text)

declare i32 @printf(i8* %format, ...)

declare i8* @malloc(i32 %size)

declare void @free(i8* %address)

declare i32 @memcmp(i8* %dest, i8* %source, i32 %size)

declare void @memcpy(i8* %dest, i8* %source, i32 %size)

declare void @memset(i8* %source, i32 %value, i32 %size)

define i8* @global.counter.create() {
entry:
	%0 = alloca i8*
	%1 = getelementptr %global.counter, %global.counter* null, i32 1
	%2 = ptrtoint %global.counter* %1 to i32
	%3 = call i8* @malloc(i32 %2)
	call void @memset(i8* %3, i32 0, i32 %2)
	%4 = bitcast i8* %3 to %global.counter*
	%5 = getelementptr %global.counter, %global.counter* %4, i32 0, i32 0
	store %global.counter.vtable.type* @global.counter.vtable.data, %global.counter.vtable.type** %5
	store i8* %3, i8** %0
	br label %body


body:
	%6 = call i32 @puts(i8* bitcast ([15 x i8]* @string.5bdaebb122965539cdd6ce77f212b65e to i8*))
	br label %exit


exit:
	%7 = load i8*, i8** %0
	ret i8* %7

}

define void @global.counter.destroy(i8* %this) {
entry:
	br label %body


body:
	%0 = call i32 @puts(i8* bitcast ([16 x i8]* @string.f8f86b3941cca26e8c147322b9a8309f to i8*))
	call void @free(i8* %this)
	br label %exit


exit:
	ret void

}

define void @global.counter.retain_shared(i8* %this) {
entry:
	br label %body


body:
	%0 = bitcast i8* %this to %global.counter*
	%1 = getelementptr %global.counter, %global.counter* %0, i32 0, i32 1
	%2 = load i32, i32* %1
	%3 = add i32 %2, 1
	store i32 %3, i32* %1
	%4 = bitcast i8* %this to %global.counter*
	%5 = getelementptr %global.counter, %global.counter* %4, i32 0, i32 1
	%6 = load i32, i32* %5
	%7 = call i32 (i8*, ...) @printf(i8* bitcast ([20 x i8]* @string.c1432ab71496ebb8b3e30bbcf37605e7 to i8*), i32 %6)
	br label %exit


exit:
	ret void

}

define void @global.counter.release_shared(i8* %this) {
entry:
	br label %body


body:
	%0 = icmp eq i8* %this, null
	br i1 %0, label %14, label %1


exit:
	ret void


1:
	%2 = bitcast i8* %this to %global.counter*
	%3 = getelementptr %global.counter, %global.counter* %2, i32 0, i32 1
	%4 = load i32, i32* %3
	%5 = sub i32 %4, 1
	store i32 %5, i32* %3
	%6 = bitcast i8* %this to %global.counter*
	%7 = getelementptr %global.counter, %global.counter* %6, i32 0, i32 1
	%8 = load i32, i32* %7
	%9 = call i32 (i8*, ...) @printf(i8* bitcast ([21 x i8]* @string.ce99a7174da84f1767b9d5235ea3c24f to i8*), i32 %8)
	%10 = bitcast i8* %this to %global.counter*
	%11 = getelementptr %global.counter, %global.counter* %10, i32 0, i32 1
	%12 = load i32, i32* %11
	%13 = icmp eq i32 %12, 0
	br i1 %13, label %16, label %15


14:
	br label %exit


15:
	br label %exit


16:
	%17 = bitcast i8* %this to %global.counter*
	%18 = getelementptr %global.counter, %global.counter* %17, i32 0, i32 4
	%19 = load void (i8*)*, void (i8*)** %18
	%20 = bitcast i8* %this to %global.counter*
	%21 = getelementptr %global.counter, %global.counter* %20, i32 0, i32 3
	%22 = load i8*, i8** %21
	call void %19(i8* %22)
	%23 = bitcast i8* %this to %global.counter*
	%24 = getelementptr %global.counter, %global.counter* %23, i32 0, i32 3
	%25 = load i8*, i8** %24
	%26 = call i32 (i8*, ...) @printf(i8* bitcast ([17 x i8]* @string.21c67ac9191c65481dbab306227b4840 to i8*), i8* %25)
	%27 = bitcast i8* %this to %global.counter*
	%28 = getelementptr %global.counter, %global.counter* %27, i32 0, i32 3
	%29 = load i8*, i8** %28
	call void @free(i8* %29)
	%30 = bitcast i8* %this to %global.counter*
	%31 = getelementptr %global.counter, %global.counter* %30, i32 0, i32 3
	%32 = load i8*, i8** %31
	%33 = bitcast i8* %this to %global.counter*
	%34 = getelementptr %global.counter, %global.counter* %33, i32 0, i32 3
	store i8* null, i8** %34
	%35 = bitcast i8* %this to %global.counter*
	%36 = getelementptr %global.counter, %global.counter* %35, i32 0, i32 2
	%37 = load i32, i32* %36
	%38 = icmp eq i32 %37, 0
	br i1 %38, label %40, label %39


39:
	br label %15


40:
	%41 = call i32 (i8*, ...) @printf(i8* bitcast ([18 x i8]* @string.4fc1bf1a9ddd2be568f08ffc8ed6b9f0 to i8*), i8* %this)
	call void @free(i8* %this)
	br label %39

}

define void @global.counter.retain_weak(i8* %this) {
entry:
	br label %body


body:
	%0 = bitcast i8* %this to %global.counter*
	%1 = getelementptr %global.counter, %global.counter* %0, i32 0, i32 2
	%2 = load i32, i32* %1
	%3 = add i32 %2, 1
	store i32 %3, i32* %1
	%4 = bitcast i8* %this to %global.counter*
	%5 = getelementptr %global.counter, %global.counter* %4, i32 0, i32 2
	%6 = load i32, i32* %5
	%7 = call i32 (i8*, ...) @printf(i8* bitcast ([18 x i8]* @string.8d9c52192bdfa908703a004b070ff63e to i8*), i32 %6)
	br label %exit


exit:
	ret void

}

define void @global.counter.release_weak(i8* %this) {
entry:
	br label %body


body:
	%0 = icmp eq i8* %this, null
	br i1 %0, label %19, label %1


exit:
	ret void


1:
	%2 = bitcast i8* %this to %global.counter*
	%3 = getelementptr %global.counter, %global.counter* %2, i32 0, i32 2
	%4 = load i32, i32* %3
	%5 = sub i32 %4, 1
	store i32 %5, i32* %3
	%6 = bitcast i8* %this to %global.counter*
	%7 = getelementptr %global.counter, %global.counter* %6, i32 0, i32 2
	%8 = load i32, i32* %7
	%9 = call i32 (i8*, ...) @printf(i8* bitcast ([19 x i8]* @string.ad2ae91c8542c824c6842efed6523f49 to i8*), i32 %8)
	%10 = bitcast i8* %this to %global.counter*
	%11 = getelementptr %global.counter, %global.counter* %10, i32 0, i32 1
	%12 = load i32, i32* %11
	%13 = icmp eq i32 %12, 0
	%14 = bitcast i8* %this to %global.counter*
	%15 = getelementptr %global.counter, %global.counter* %14, i32 0, i32 2
	%16 = load i32, i32* %15
	%17 = icmp eq i32 %16, 0
	%18 = and i1 %13, %17
	br i1 %18, label %21, label %20


19:
	br label %exit


20:
	br label %exit


21:
	%22 = call i32 (i8*, ...) @printf(i8* bitcast ([18 x i8]* @string.4fc1bf1a9ddd2be568f08ffc8ed6b9f0 to i8*), i8* %this)
	call void @free(i8* %this)
	br label %20

}

define i32 @main() {
entry:
	%0 = alloca i32
	%1 = alloca i8*
	br label %body


body:
	%2 = call i8* @global.derive_class.create()
	%3 = call i8* @global.counter.create()
	call void @global.counter.retain_shared(i8* %3)
	%4 = bitcast i8* %3 to %global.counter*
	%5 = getelementptr %global.counter, %global.counter* %4, i32 0, i32 3
	store i8* %2, i8** %5
	%6 = bitcast i8* %3 to %global.counter*
	%7 = getelementptr %global.counter, %global.counter* %6, i32 0, i32 4
	store void (i8*)* @global.derive_class.destroy, void (i8*)** %7
	store i8* %3, i8** %1
	%8 = load i8*, i8** %1
	call void @global.echo(i8* %8)
	store i32 0, i32* %0
	br label %exit


exit:
	%9 = load i8*, i8** %1
	call void @global.counter.release_shared(i8* %9)
	%10 = load i32, i32* %0
	ret i32 %10

}

define void @global.echo(i8* %dc) {
entry:
	br label %body


body:
	%0 = bitcast i8* %dc to %global.counter*
	%1 = getelementptr %global.counter, %global.counter* %0, i32 0, i32 3
	%2 = load i8*, i8** %1
	%3 = bitcast i8* %2 to %global.derive_class*
	%4 = getelementptr %global.derive_class, %global.derive_class* %3, i32 0, i32 0
	%5 = load %global.derive_class.vtable.type*, %global.derive_class.vtable.type** %4
	%6 = getelementptr %global.derive_class.vtable.type, %global.derive_class.vtable.type* %5, i32 0, i32 2
	%7 = load void (i8*)*, void (i8*)** %6
	call void %7(i8* %2)
	br label %exit


exit:
	ret void

}

define i8* @global.base_class.create() {
entry:
	%0 = alloca i8*
	%1 = getelementptr %global.base_class, %global.base_class* null, i32 1
	%2 = ptrtoint %global.base_class* %1 to i32
	%3 = call i8* @malloc(i32 %2)
	call void @memset(i8* %3, i32 0, i32 %2)
	%4 = bitcast i8* %3 to %global.base_class*
	%5 = getelementptr %global.base_class, %global.base_class* %4, i32 0, i32 0
	store %global.base_class.vtable.type* @global.base_class.vtable.data, %global.base_class.vtable.type** %5
	store i8* %3, i8** %0
	br label %body


body:
	%6 = call i32 @puts(i8* bitcast ([18 x i8]* @string.726bd3560bd4c136648f7760895d8d62 to i8*))
	br label %exit


exit:
	%7 = load i8*, i8** %0
	ret i8* %7

}

define void @global.base_class.destroy(i8* %this) {
entry:
	br label %body


body:
	%0 = call i32 @puts(i8* bitcast ([17 x i8]* @string.362aeeddb3d01da539cb6755bde46953 to i8*))
	br label %exit


exit:
	ret void

}

define void @global.base_class.echo(i8* %this) {
entry:
	br label %body


body:
	%0 = call i32 @puts(i8* bitcast ([13 x i8]* @string.9bcbb503bda6c8ad83f772846b706f08 to i8*))
	br label %exit


exit:
	ret void

}

define i8* @global.derive_class.create() {
entry:
	%0 = alloca i8*
	%1 = getelementptr %global.derive_class, %global.derive_class* null, i32 1
	%2 = ptrtoint %global.derive_class* %1 to i32
	%3 = call i8* @malloc(i32 %2)
	call void @memset(i8* %3, i32 0, i32 %2)
	%4 = bitcast i8* %3 to %global.derive_class*
	%5 = getelementptr %global.derive_class, %global.derive_class* %4, i32 0, i32 0
	store %global.derive_class.vtable.type* @global.derive_class.vtable.data, %global.derive_class.vtable.type** %5
	store i8* %3, i8** %0
	br label %body


body:
	%6 = call i32 @puts(i8* bitcast ([20 x i8]* @string.33b7808bf372c3d58730520160cb2c15 to i8*))
	br label %exit


exit:
	%7 = load i8*, i8** %0
	ret i8* %7

}

define void @global.derive_class.destroy(i8* %this) {
entry:
	br label %body


body:
	call void @global.base_class.destroy(i8* %this)
	%0 = call i32 @puts(i8* bitcast ([19 x i8]* @string.ef25b0542457581e67c27a0dddb7bda5 to i8*))
	br label %exit


exit:
	ret void

}

define void @global.derive_class.echo(i8* %this) {
entry:
	br label %body


body:
	call void @global.base_class.echo(i8* %this)
	%0 = call i32 @puts(i8* bitcast ([15 x i8]* @string.895758554639f423e017c6610cbf460b to i8*))
	br label %exit


exit:
	ret void

}

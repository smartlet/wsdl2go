# package wsdl2go

使用"github.com/hooklift/gowsdl"改造生成ews相关的实现! 但要通用还需要继续优化.

- main.go
    - 修改"-o"会自动拼加"-p"
    - 修改"serverFile"的命名
    - 移除"-make-public"所有都固定是true.
    - 添加"-require"生成
      - 标准空间"http://www.w3.org/2001/XMLSchema"的数据类型
      - 服务调用依赖SOAPClient接口

- gowsdl/
    - 搜索并修改"xsd2GoTypes"映射类型, 将"soap.XSD..."都去掉.
    - 搜索并修改"soap.xxx"相关的模板. 
        - 把"soap.XSD..."前缀去掉,匹配source新加的类型
        - 把"soap.Client"改成"SOAPClient"接口

- source/
    - 添加对应"http://www.w3.org/2001/XMLSchema"标准空间缺失的类型.
    